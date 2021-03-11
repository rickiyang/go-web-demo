package dao

import (
	"go.uber.org/zap"
	"gorm-demo/constants"
	"gorm-demo/model"
	"gorm.io/gorm/clause"
)

//gorm 操作使用

//单行  新增用户
func InsertOneUser(user model.User) (id int64, err error) {
	tx := constants.GVA_DB.Create(&user)
	if tx.Error != nil {
		constants.GVA_LOG.Error("InsertOne err", zap.Any("err", tx.Error))
		return 0, tx.Error
	}
	return user.Id, nil
}

//插入主键冲突的时候操作
func UpsertOp(user model.User) {
	// 在冲突时，什么都不做
	constants.GVA_DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

	// 在`id`冲突时，将列更新为默认值
	constants.GVA_DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"name": "", "age": 0, "sex": 1}),
	}).Create(&user)

	// 在`id`冲突时，将列更新为新值
	constants.GVA_DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "age", "sex", "phone"}),
	}).Create(&user)

	// 在冲突时，更新除主键以外的所有列到新值。
	constants.GVA_DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user)

}

//批量插入
func BatchInsertuser(user []model.User) (ids []int64, err error) {
	tx := constants.GVA_DB.CreateInBatches(user, len(user))
	if tx.Error != nil {
		constants.GVA_LOG.Error("BatchInsert err", zap.Any("err", tx.Error))
		return []int64{}, tx.Error
	}
	ids = []int64{}
	for idx, user := range user {
		ids[idx] = user.Id
	}
	return ids, nil
}

//根据 id 删除数据
func DeleteUserById(id int64) (err error) {
	user := model.User{Id: id}
	err = constants.GVA_DB.Delete(&user).Error
	if err != nil {
		constants.GVA_LOG.Error("DeleteUserById err", zap.Any("err", err))
		return err
	}
	return nil
}

//根据 id 批量删除数据
func BatchDeleteUserByIds(ids []int64) (err error) {
	if ids == nil || len(ids) == 0 {
		return
	}
	//删除方式1
	err = constants.GVA_DB.Where("id in ?", ids).Delete(model.User{}).Error
	if err != nil {
		constants.GVA_LOG.Error("DeleteUserById err", zap.Any("err", err))
		return err
	}

	//删除方式 2
	//constants.GVA_DB.Delete(model.User{}, "id in ?", ids)

	return nil
}

//根据id更新数据，全量字段更新，即使字段是0值
func UpdateUserById(user model.User) (err error) {
	err = constants.GVA_DB.Save(&user).Error
	if err != nil {
		constants.GVA_LOG.Error("UpdateUserById err", zap.Any("err", err))
		return err
	}
	return nil
}

//更新指定列
//update user set `columnName` = v where id = id;
func UpdateSpecialColumn(id int64, columnName string, v interface{}) (err error) {
	err = constants.GVA_DB.Model(&model.User{Id: id}).Update(columnName, v).Error
	if err != nil {
		constants.GVA_LOG.Error("UpdateSpecialColumn err", zap.Any("err", err))
		return err
	}
	return nil
}

//更新- 根据 `struct` 更新属性，只会更新非零值的字段
//update user set `columnName` = v where id = id;
//当通过 struct 更新时，GORM 只会更新非零字段。 如果您想确保指定字段被更新，你应该使用 Select 更新选定字段，或使用 map 来完成更新操作
func UpdateSelective(user model.User) (effected int64, err error) {
	tx := constants.GVA_DB.Model(&user).Updates(&model.User{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Sex:   user.Sex,
		Phone: user.Phone,
	})

	//map 方式会更新0值字段
	tx = constants.GVA_DB.Model(&user).Updates(map[string]interface{}{
		"Id":    user.Id,
		"Name":  user.Name,
		"Age":   user.Age,
		"Sex":   user.Sex,
		"Phone": user.Phone,
	})

	//Select 方式指定列名
	tx = constants.GVA_DB.Model(&user).Select("Name", "Age", "Phone").Updates(&model.User{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Sex:   user.Sex,
		Phone: user.Phone,
	})

	// Select 所有字段（查询包括零值字段的所有字段）
	tx = constants.GVA_DB.Model(&user).Select("*").Updates(&model.User{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Sex:   user.Sex,
		Phone: user.Phone,
	})

	// Select 除 Phone 外的所有字段（包括零值字段的所有字段）
	tx = constants.GVA_DB.Model(&user).Select("*").Omit("Phone").Updates(&model.User{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Sex:   user.Sex,
		Phone: user.Phone,
	})
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

//根据 条件 批量更新
func BatchUpdateByIds(ids []int64, user model.User) (effected int64, err error) {
	if ids == nil || len(ids) == 0 {
		return
	}
	tx := constants.GVA_DB.Model(model.User{}).Where("id in ?", ids).Updates(&user)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

//查询操作
func queryOp(user model.User) {

	// 获取第一条记录（主键升序）
	// SELECT * FROM user ORDER BY id LIMIT 1;
	constants.GVA_DB.First(&user)

	// 获取一条记录，没有指定排序字段
	// SELECT * FROM user LIMIT 1;
	constants.GVA_DB.Take(&user)

	// 获取最后一条记录（主键降序）
	// SELECT * FROM user ORDER BY id DESC LIMIT 1;
	constants.GVA_DB.Last(&user)

	// SELECT * FROM user WHERE id = 10;
	constants.GVA_DB.First(&user, 10)

	// SELECT * FROM user WHERE id = 10;
	constants.GVA_DB.First(&user, "10")

	// SELECT * FROM user WHERE id IN (1,2,3);
	constants.GVA_DB.Find(&user, []int{1, 2, 3})

	// 获取全部记录
	// SELECT * FROM user;
	result := constants.GVA_DB.Find(&user)
	result.Rows()

	// 获取第一条匹配的记录
	// SELECT * FROM user WHERE name = 'xiaoming' ORDER BY id LIMIT 1;
	constants.GVA_DB.Where("name = ?", "xiaoming").First(&user)

	// 获取全部匹配的记录
	// SELECT * FROM user WHERE name <> 'xiaoming';
	constants.GVA_DB.Where("name <> ?", "xiaoming").Find(&user)

	// IN
	// SELECT * FROM user WHERE name IN ('xiaoming','xiaohong');
	constants.GVA_DB.Where("name IN ?", []string{"xiaoming", "xiaohong"}).Find(&user)

	// LIKE
	// SELECT * FROM user WHERE name LIKE '%ming%';
	constants.GVA_DB.Where("name LIKE ?", "%ming%").Find(&user)

	// AND
	// SELECT * FROM user WHERE name = 'xiaoming' AND age >= 33;
	constants.GVA_DB.Where("name = ? AND age >= ?", "xiaoming", 33).Find(&user)

	// Time
	// SELECT * FROM user WHERE updated_at > '2021-03-10 15:44:23';
	constants.GVA_DB.Where("updated_at > ?", "2021-03-10 15:44:23").Find(&user)

	// BETWEEN
	// SELECT * FROM user WHERE created_at BETWEEN ''2021-03-07 15:44:23' AND '2021-03-10 15:44:23';
	constants.GVA_DB.Where("created_at BETWEEN ? AND ?", "2021-03-07 15:44:23", "2021-03-10 15:44:23").Find(&user)

	// SELECT * FROM user WHERE NOT name = "xiaoming" ORDER BY id LIMIT 1;
	constants.GVA_DB.Not("name = ?", "xiaoming").First(&user)

	// Not In
	// SELECT * FROM user WHERE name NOT IN ("xiaoming", "xiaohong");
	constants.GVA_DB.Not(map[string]interface{}{"name": []string{"xiaoming", "xiaohong"}}).Find(&user)

	// Struct
	// SELECT * FROM user WHERE name <> "xiaoming" AND age <> 20 ORDER BY id LIMIT 1;
	constants.GVA_DB.Not(model.User{Name: "xiaoming", Age: 20}).First(&user)

	// 不在主键切片中的记录
	// SELECT * FROM user WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
	constants.GVA_DB.Not([]int64{1, 2, 3}).First(&user)

	// SELECT * FROM user WHERE name = 'xiaoming' OR name = 'xiaohong';
	constants.GVA_DB.Where("name = ?", "xiaoming").Or("name = ?", "xiaohong").Find(&user)

	// Struct
	// SELECT * FROM user WHERE name = 'xiaoming' OR (name = 'xiaohong' AND age = 20);
	constants.GVA_DB.Where("name = 'xiaoming'").Or(model.User{Name: "xiaohong", Age: 20}).Find(&user)

	// Map
	// SELECT * FROM user WHERE name = 'xiaoming' OR (name = 'xiaohong' AND age = 20);
	constants.GVA_DB.Where("name = 'xiaoming'").Or(map[string]interface{}{"name": "xiaohong", "age": 20}).Find(&user)

	// SELECT name, age FROM user;
	constants.GVA_DB.Select("name", "age").Find(&user)

	// SELECT name, age FROM user;
	constants.GVA_DB.Select([]string{"name", "age"}).Find(&user)

	// SELECT COALESCE(age,'20') FROM user;
	constants.GVA_DB.Table("user").Select("COALESCE(age,?)", 20).Rows()

	// SELECT * FROM user ORDER BY age desc, name;
	constants.GVA_DB.Order("age desc, name").Find(&user)

	// 多个 order
	// SELECT * FROM user ORDER BY age desc, name;
	constants.GVA_DB.Order("age desc").Order("name").Find(&user)

	// SELECT * FROM user ORDER BY FIELD(id,1,2,3)
	constants.GVA_DB.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{[]int{1, 2, 3}}, WithoutParentheses: true},
	}).Find(&model.User{})

	// SELECT * FROM user LIMIT 10;
	constants.GVA_DB.Limit(10).Find(&user)

	// SELECT * FROM user OFFSET 10;
	constants.GVA_DB.Offset(10).Find(&user)

	// SELECT * FROM user OFFSET 0 LIMIT 10;
	constants.GVA_DB.Limit(10).Offset(0).Find(&user)

	// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "ming%" GROUP BY `name`
	constants.GVA_DB.Model(&model.User{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)

	// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"
	constants.GVA_DB.Model(&model.User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "group").Find(&result)

	//SELECT distinct(name, age) from user order by name, age desc
	constants.GVA_DB.Distinct("name", "age").Order("name, age desc").Find(&user)

}

//事务测试
func TestGormTx(user model.User) (err error) {
	tx := constants.GVA_DB.Begin()
	// 注意，一旦你在一个事务中，使用tx作为数据库句柄
	if err := tx.Create(&model.User{
		Name:  "liliya",
		Age:   13,
		Sex:   0,
		Phone: "15543212346",
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Updates(&model.User{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Sex:   user.Sex,
		Phone: user.Phone,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
