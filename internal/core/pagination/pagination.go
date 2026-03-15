package pagination

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type QueryParams struct {
	Page       int                    `json:"page" form:"page"`
	Size       int                    `json:"size" form:"size"`
	SortBy     string                 `json:"sortBy" form:"sortBy"`         // 排序的字段
	SortOrder  string                 `json:"sortOrder" form:"sortOrder"`   // "asc" 或 "desc"
	Preloads   []string               `json:"preloads" form:"preloads"`     // 预加载
	Filters    map[string]interface{} `json:"filters" form:"filters"`       // 用于存放 Where 条件
	Keyword    string                 `json:"keyword" form:"keyword"`       // 新增：模糊搜索关键词
	Conditions []Condition            `json:"conditions" form:"conditions"` // 动态条件
	Joins      []JoinParam            `json:"joins" form:"joins"`           // 新增关联查询参数
}

type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // 条件: =, !=, >, >=, <, <=, IN, LIKE, BETWEEN
	Value    interface{} `json:"value"`
}

type JoinParam struct {
	Table     string        `json:"table"`     // 要关联的表名
	Condition string        `json:"condition"` // 关联条件
	Args      []interface{} `json:"args"`      // 条件参数
	Type      string        `json:"type"`      // JOIN类型：INNER/LEFT/RIGHT
}

// Pagination 分页结构体
type Pagination struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

// 应用全局搜索
func applyGlobalSearch(query *gorm.DB, keyword string, fields []string) *gorm.DB {
	if keyword == "" || len(fields) == 0 {
		return query
	}
	// 构建OR条件
	var orConditions []string
	var args []interface{}
	for _, field := range fields {
		orConditions = append(orConditions, quoteField(field)+" LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	return query.Where(strings.Join(orConditions, " OR "), args...)
}

func quoteField(field string) string {
	if field == "" {
		return field
	}
	parts := strings.Split(field, ".")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = "`" + part + "`"
	}
	return strings.Join(parts, ".")
}

// 应用排序
func applySorting(query *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	if sortBy == "" {
		return query
	}

	order := sortBy
	if sortOrder != "" {
		order += " " + sortOrder
	}
	return query.Order(order)
}

// 应用动态条件
func applyConditions(query *gorm.DB, conditions []Condition) *gorm.DB {
	for _, cond := range conditions {
		switch cond.Operator {
		case "=":
			query = query.Where(cond.Field+" = ?", cond.Value)
		case "!=":
			query = query.Where(cond.Field+" != ?", cond.Value)
		case ">":
			query = query.Where(cond.Field+" > ?", cond.Value)
		case ">=":
			query = query.Where(cond.Field+" >= ?", cond.Value)
		case "<":
			query = query.Where(cond.Field+" < ?", cond.Value)
		case "<=":
			query = query.Where(cond.Field+" <= ?", cond.Value)
		case "IN":
			query = query.Where(cond.Field+" IN (?)", cond.Value)
		case "LIKE":
			query = query.Where(cond.Field+" LIKE ?", cond.Value)
		case "IS_NULL":
			query = query.Where(cond.Field + " IS NULL")
		case "IS_NOT_NULL":
			query = query.Where(cond.Field + " IS NOT NULL")
		case "BETWEEN":
			if values, ok := cond.Value.([]interface{}); ok && len(values) == 2 {
				query = query.Where(cond.Field+" BETWEEN ? AND ?", values[0], values[1])
			}
		}
	}
	return query
}

func applyJoins(query *gorm.DB, joins []JoinParam) *gorm.DB {
	for _, join := range joins {
		joinType := strings.ToUpper(join.Type)
		joinClause := fmt.Sprintf("%s JOIN %s ON %s", joinType, join.Table, join.Condition)

		// 如果有参数，使用带参数的 Joins 方法
		if len(join.Args) > 0 {
			query = query.Joins(joinClause, join.Args...)
		} else {
			query = query.Joins(joinClause)
		}
	}
	return query
}

// Paginate 分页函数
func Paginate(tx *gorm.DB, data interface{}, params QueryParams) (Pagination, error) {
	var pagination Pagination
	var total int64

	// ---------- 1. 默认分页参数 ----------
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 10
	}

	query := tx
	// 兜底：只有没指定 Table 才用 Model
	if query.Statement.Table == "" {
		query = query.Model(data)
	}

	// ---------- 2. 强制解析模型（用于全局搜索字段） ----------
	if err := query.Statement.Parse(data); err != nil {
		return pagination, fmt.Errorf("failed to parse model, err: %w", err)
	}

	// ---------- 3. 全局搜索 ----------
	if params.Keyword != "" {
		var searchFields []string
		stmt := query.Statement
		if stmt.Schema != nil {
			for _, field := range stmt.Schema.Fields {
				if field.DBName != "" {
					searchFields = append(searchFields, field.DBName)
				}
			}
		}
		if len(searchFields) > 0 {
			query = applyGlobalSearch(query, params.Keyword, searchFields)
		}
	}

	// ---------- 4. JOIN（⚠️ 放在 where 之前） ----------
	if len(params.Joins) > 0 {
		query = applyJoins(query, params.Joins)
	}

	// ---------- 5. Filters（仅等值匹配，保持旧行为） ----------
	for column, value := range params.Filters {
		query = query.Where(column+" = ?", value)
	}

	// ---------- 6. 动态 Conditions（新增能力） ----------
	if len(params.Conditions) > 0 {
		query = applyConditions(query, params.Conditions)
	}

	// ---------- 7. COUNT（防 JOIN 放大，可选） ----------
	countQuery := query

	// ⚠️ 仅当存在 JOIN 时才 Distinct，避免影响旧模块
	if len(params.Joins) > 0 {
		// 默认用主表 id 去重（假设所有表都有 id）
		if query.Statement.Table != "" {
			countQuery = countQuery.Distinct(query.Statement.Table + ".id")
		}
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return pagination, err
	}

	// ---------- 8. 排序 ----------
	query = applySorting(query, params.SortBy, params.SortOrder)

	// ---------- 9. Preload ----------
	for _, preload := range params.Preloads {
		query = query.Preload(preload)
	}

	// ---------- 10. 分页 ----------
	offset := (params.Page - 1) * params.Size
	if err := query.Offset(offset).Limit(params.Size).Find(data).Error; err != nil {
		return pagination, err
	}

	return Pagination{
		Page:  params.Page,
		Size:  params.Size,
		Total: total,
	}, nil
}
