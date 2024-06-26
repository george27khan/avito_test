package banner_content_hist

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type BannerContentHist struct {
	id        int         `json:"id"`
	BannerId  int64       `json:"banner_id"`
	Content   interface{} `json:"content"`
	Version   int         `json:"varsion"`
	CreatedAt time.Time   `json:"created_at"`
}

// initVersion функция для определения версии контента
func (c *BannerContentHist) initVersion(ctx context.Context, tx pgx.Tx) error {
	query := "select count(1)+1 from avito_banner.banner_content_hist where banner_id = $1"
	if err := tx.QueryRow(ctx, query, c.BannerId).Scan(&c.Version); err != nil {
		return fmt.Errorf("initVersion ошибка оперделения версии баннера (%v)", err.Error())
	}
	return nil
}

// Insert функция для создания контента баннера в версионировании
func (c *BannerContentHist) Insert(ctx context.Context, tx pgx.Tx) error {
	if err := c.initVersion(ctx, tx); err != nil {
		return fmt.Errorf("в Insert (%v)", err.Error())
	}
	query := "INSERT INTO avito_banner.banner_content_hist(banner_id, content, version) VALUES ($1, $2, $3) "

	if _, err := tx.Exec(ctx, query, c.BannerId, c.Content, c.Version); err != nil {
		return fmt.Errorf("в Insert ошибка создания записи баннера в версионировании (%v)", err.Error())
	}
	return nil
}

// DeleteByBannerId функция для удаления контента баннера
func DeleteByBannerId(ctx context.Context, tx pgx.Tx, bannerID int64) error {
	query := "delete from avito_banner.banner_content_hist where banner_id = $1"
	if _, err := tx.Exec(ctx, query, bannerID); err != nil {
		return fmt.Errorf("DeleteByBannerId: ошибка в процессе удаления контента (%v)", err.Error())
	}
	return nil
}
