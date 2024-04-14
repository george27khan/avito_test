package tag_feature

import (
	db "avito_test/pkg/postgres_db/connection"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TagFeauture struct {
	TagFeatureId int64
	TagId        int64
	FeatureId    int64
	IsActive     bool
	CreatedDt    time.Time
}

// Insert функция для добавление записи в таблицу тегов
func Insert(ctx context.Context, tx pgx.Tx, tags []int64, feature int64) error {
	query := "INSERT INTO avito_banner.tag_feature(tag_id, feature_id) VALUES ($1, $2)"
	for _, tag := range tags {
		if _, err := tx.Exec(ctx, query, tag, feature); err != nil {
			return fmt.Errorf("в Insert ошибка создания записи тега (%v)", err.Error())
		}
	}
	return nil
}

// Get функция для формирования объекта TagFeauture
func Get(ctx context.Context, conn *pgxpool.Conn, tagId int64, featureId int64) (TagFeauture, error) {
	var tf TagFeauture
	query := "select tag_feature_id, created_at, is_active" +
		"from avito_banner.tag_feature tf " +
		"where tf.tag_id = $1 and tf.feature_id = $2"
	if err := conn.QueryRow(ctx, query, tagId, featureId).
		Scan(&tf.TagFeatureId, &tf.CreatedDt, &tf.IsActive); err != nil {
		return tf, err
	}
	return tf, nil
}

// DeleteByBannerId функция удаления тегов баннера
func DeleteByBannerId(ctx context.Context, tx pgx.Tx, bannerId int64) error {
	query := "delete from avito_banner.tag_feature " +
		" where tag_feature_id in (select tf.tag_feature_id " +
		"from avito_banner.banner b " +
		"join avito_banner.tag_feature tf on tf.feature_id = b.feature_id " +
		"where b.banner_id = $1)"
	if _, err := tx.Exec(ctx, query, bannerId); err != nil {
		return fmt.Errorf("DeleteByBannerId: ошибка в процессе удаления тэга (%v)", err.Error())
	}
	return nil
}

// Check функция проверки существования связки тег-фича в базе
func (tf *TagFeauture) Check(ctx context.Context, conn *pgxpool.Conn) error {
	var cnt int
	query := "select 1 from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2"
	if err := conn.QueryRow(ctx, query, tf.TagId, tf.FeatureId).Scan(&cnt); err != nil {
		return fmt.Errorf("баннер для тега %v и фичи %v уже определен %v", tf.TagId, tf.FeatureId, err.Error())
	}
	return nil
}

// Activate функция для активации/деактивации тэга
func Activate(ctx context.Context, tx pgx.Tx, tagFeatureId int64, isActive bool) error {
	query := "update avito_banner.tag_feature " +
		"set is_active=$1 " +
		"where tag_feature_id = $2"
	if _, err := tx.Exec(ctx, query, isActive, tagFeatureId); err != nil {
		return fmt.Errorf("ошибка в процессе обновления тэга %v", err.Error())
	}
	return nil
}

// CheckData функция проверки существования связки тег-фича в базе по списку тегов
func CheckData(ctx context.Context, tags []int64, feature int64) error {
	conn, err := db.PGPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, tag := range tags {
		tagFeature := TagFeauture{TagId: tag, FeatureId: feature}
		if err := tagFeature.Check(ctx, conn); err == nil {
			return fmt.Errorf("ошибка проверки существования тега %v", err.Error())
		}
	}
	return nil
}

// MergeTags функция обновления данных банера о тегах
func MergeTags(ctx context.Context, tx pgx.Tx, tagIds []int64, feature int64, bannerId int64) error {
	var tag TagFeauture
	tagIdsMap := make(map[int64]bool)
	tagFeatureActive := make(map[int64]bool)
	for _, tag := range tagIds {
		tagIdsMap[tag] = false //флаг присутсвия тега в базе
	}
	// поиск активных связей баннера
	query := "select tf.tag_feature_id, tf.tag_id, tf.feature_id, tf.is_active " +
		"from avito_banner.banner b " +
		"join avito_banner.tag_feature tf on tf.feature_id = b.feature_id " +
		"where b.banner_id = $1"
	rows, err := tx.Query(ctx, query, bannerId)
	if err != nil {
		return fmt.Errorf("в MergeTags ошибка проверки существования тега (%v)", err.Error())
	}
	defer rows.Close()

	//тэги на отвязку
	for rows.Next() {
		if err = rows.Scan(&tag.TagFeatureId, &tag.TagId, &tag.FeatureId, &tag.IsActive); err != nil {
			return fmt.Errorf("в MergeTags ошибка чтения данных из таблицы tag_feature %v", err.Error())
		}
		//если тэг из нового списка не активен, то активируем
		if _, ok := tagIdsMap[tag.TagId]; ok {
			//активируем
			tagIdsMap[tag.TagId] = true //делаем метку о наличии
			if !tag.IsActive {
				tagFeatureActive[tag.TagFeatureId] = true
			}
		} else { //если тега в списке нет, то деактивируем
			//деактивируем
			tagFeatureActive[tag.TagFeatureId] = false
		}
	}

	//актуализируем активность тегов
	for tagFeatureId, isActive := range tagFeatureActive {
		fmt.Println(tagFeatureId, isActive)
		if err := Activate(ctx, tx, tagFeatureId, isActive); err != nil {
			return fmt.Errorf("в MergeTags ошибка при актуализации тега %v", err.Error())
		}
	}
	//добавление новых тэгов
	for tag, isExist := range tagIdsMap {
		if !isExist {
			if err := Insert(ctx, tx, []int64{tag}, feature); err != nil {
				return fmt.Errorf("в MergeTags ошибка при создании тега %v", err.Error())
			}
		}
	}
	return nil
}
