package banner_handler

import (
	"avito_test/pkg/auth"
	db "avito_test/pkg/db_avito_banner"
	bn "avito_test/pkg/db_avito_banner/banner"
	bch "avito_test/pkg/db_avito_banner/banner_content_hist"
	tf "avito_test/pkg/db_avito_banner/tag_feature"
	usr "avito_test/pkg/db_avito_banner/user"
	"avito_test/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const cashExpiration = time.Minute * 5

type PostBody struct {
	TagIds    []int64     `json:"tag_ids"`
	FeatureId int64       `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
}

// BannerLink структура для PatchBanner чтобы определять null значения
type BannerLink struct {
	TagIds    []int64 `json:"tag_ids"`
	FeatureId *int64  `json:"feature_id"`
	Content   *string `json:"content"`
	IsActive  *bool   `json:"is_active"`
}

// ContentCash структура для хранения кешированных данных
type ContentCash struct {
	content  string `json:"content"`
	isActive bool   `json:"is_active"`
}

func toInt(val string) (int64, error) {
	res, err := strconv.ParseInt(val, 10, 64)
	return res, err
}

// GetBannerVersion добавление баннера
func GetBannerVersion(c *gin.Context) {
	var useLastVersionBool bool
	ctx := context.Background()
	conn, err := db.ConnectPool(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка подключения к базе - %v", err.Error())})
		return
	}
	defer conn.Release()

	//чтение параметров
	feature := c.Query("feature_id")
	tag := c.Query("tag_id")
	useLastVersion := c.Query("use_last_revision")
	//token := c.GetHeader("token")

	//валидация параметров
	if tag == "" || feature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не переданны обязательные параметры для запроса"})
		return
	}
	tagInt, err := toInt(tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный тип параметра tag_id"})
		return
	}
	featureInt, err := toInt(feature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный тип параметра feature_id"})
		return
	}
	if useLastVersion != "" {
		useLastVersionBool, err = strconv.ParseBool(useLastVersion)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный тип параметра use_last_revision"})
			return
		}
	} else {
		useLastVersionBool = false
	}

	//попытка чтения из кеша
	if !useLastVersionBool {
		if content := redis.Get(ctx, tag+"_"+feature+"_all"); content != "" {
			fmt.Println("cash")
			c.JSON(http.StatusOK, content)
		}
	}
	//чтение из базы
	content, isActive, err := bn.GetContentByTagFeature(ctx, conn, tagInt, featureInt)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	//нет доступа
	if !isActive {
		c.Status(http.StatusForbidden)
		return
	}
	//кэшируем результат
	redis.Set(ctx, tag+"_"+feature+"_all", content, cashExpiration)
	c.JSON(http.StatusOK, content)
}

// GetUserBanner получение баннера
func GetUserBanner(c *gin.Context) {
	var (
		useLastVerisionBool bool
		contentCash         ContentCash
	)

	ctx := context.Background()
	conn, err := db.ConnectPool(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ошибка подключения к базе - %v", err.Error())})
		return
	}
	defer conn.Release()

	//чтение параметров
	isAdmin := c.MustGet("is_admin").(bool)
	feature := c.Query("feature_id")
	tag := c.Query("tag_id")
	useLastVerision := c.Query("use_last_revision")

	//валидация параметров
	if tag == "" || feature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не переданны обязательные параметры для запроса"})
		return
	}
	tagInt, err := toInt(tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный тип параметра tag_id"})
		return
	}
	featureInt, err := toInt(feature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный тип параметра feature_id"})
		return
	}
	if useLastVerision != "" {
		useLastVerisionBool, err = strconv.ParseBool(useLastVerision)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный тип параметра use_last_revision"})
			return
		}
	} else {
		useLastVerisionBool = false
	}

	//попытка чтения из кеша
	if !useLastVerisionBool {
		if contentCashJSON := redis.Get(ctx, tag+"_"+feature); contentCashJSON != "" {
			fmt.Println("cash")
			if err := json.Unmarshal([]byte(contentCashJSON), &contentCash); err == nil {
				//проверка доступа на получение данных
				if !isAdmin && !contentCash.isActive {
					c.Status(http.StatusForbidden)
					return
				}
				c.JSON(http.StatusOK, contentCash.content)
				return
			}
		}
	}

	//чтение из базы
	content, isActive, err := bn.GetContentByTagFeature(ctx, conn, tagInt, featureInt)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	//проверка доступа на получение данных
	if !isAdmin && !isActive {
		c.Status(http.StatusForbidden)
		return
	}

	//кэшируем результат
	contentCashJSON, _ := json.Marshal(ContentCash{content: content, isActive: isActive})
	redis.Set(ctx, tag+"_"+feature, string(contentCashJSON), cashExpiration)
	c.JSON(http.StatusOK, content)
}

// PostBanner добавление баннера
func PostBanner(c *gin.Context) {
	var banner bn.Banner

	//проверка доступа
	isAdmin := c.MustGet("is_admin")
	if !isAdmin.(bool) {
		c.Status(http.StatusForbidden)
		return
	}

	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ошибка разбора данных из тела запроса (%v)", err.Error())})
		return
	}
	//проверим существует ли такой баннер в базе
	if err := banner.CheckTagFeature(context.Background()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bannerId, err := banner.Insert(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"banner_id": bannerId})
}

// GetBanner запрос всех баннеров по тегу и/или фиче
func GetBanner(c *gin.Context) {
	var (
		bannerList                              []bn.Banner
		tagInt, featureInt, limitInt, offsetInt int64
		err                                     error
	)

	//проверка доступа
	isAdmin := c.MustGet("is_admin")
	if !isAdmin.(bool) {
		c.Status(http.StatusForbidden)
		return
	}

	feature := c.Query("feature_id")
	tag := c.Query("tag_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	if tag != "" && feature != "" {
		tagInt, err1 := toInt(tag)
		featureInt, err2 := toInt(feature)
		if err1 == nil && err2 == nil {
			if banner, err := bn.GetByTagFeature(context.Background(), tagInt, featureInt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				bannerList = append(bannerList, banner)
				c.JSON(http.StatusOK, bannerList)
				return
			}
		}
	} else if tag != "" {
		if tagInt, err = toInt(tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный tag_id=%v: %v", tag, err.Error())})
			return
		}
		if limitInt, err = toInt(limit); err != nil && limit != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный limit=%v: %v", limit, err.Error())})
			return
		}
		if offsetInt, err = toInt(offset); err != nil && offset != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный offset=%v: %v", offset, err.Error())})
			return
		}
		if bannerList, err := bn.GetByTag(context.Background(), tagInt, limitInt, offsetInt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, bannerList)
			return
		}
	} else if feature != "" {
		if featureInt, err = toInt(feature); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный feature_id=%v - %v", feature, err.Error())})
			return
		}
		if limitInt, err = toInt(limit); err != nil && limit != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный limit=%v - %v", limit, err.Error())})
			return
		}
		if offsetInt, err = toInt(offset); err != nil && offset != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный offset=%v - %v", offset, err.Error())})
			return
		}
		if bannerList, err := bn.GetByFeature(context.Background(), featureInt, limitInt, offsetInt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, bannerList)
			return
		}
	}
}

// PatchBanner изменение информации о баннере
func PatchBanner(c *gin.Context) {
	var bannerLink BannerLink

	//проверка доступа
	isAdmin := c.MustGet("is_admin")
	if !isAdmin.(bool) {
		c.Status(http.StatusForbidden)
		return
	}

	ctx := context.Background()
	conn, tx, err := db.ConnectPoolTrx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка подключения к БД (%v)", err.Error())})
		return
	}
	defer conn.Release()
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)

		} else {
			_ = tx.Commit(ctx)
		}
	}()

	bannerId, err := toInt(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Передан некорректный id (%v)", err.Error())})
		return
	}
	//чтение тала запроса
	if err = c.ShouldBindJSON(&bannerLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Некорректне данные в теле запроса (%v)", err.Error())})
		return
	}

	//проверка существования bannerId
	if err = bn.Exist(context.Background(), conn, bannerId); err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	banner, err := bn.Get(ctx, conn, bannerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//проверка существования фичи для баннера
	if bannerLink.FeatureId != nil && *bannerLink.FeatureId != banner.FeatureId {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Для баннера %v не определена фича %v", bannerId, *bannerLink.FeatureId)})
		return
	}

	//обновляем поля баннера
	if bannerLink.IsActive != nil && banner.IsActive != *bannerLink.IsActive {
		if err = banner.UpdateField(ctx, tx, "is_active", *bannerLink.IsActive); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if bannerLink.Content != nil && banner.Content != *bannerLink.Content {
		if err = banner.UpdateField(ctx, tx, "content", bannerLink.Content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	//актуализируем теги
	if bannerLink.FeatureId != nil && bannerLink.TagIds != nil {
		if err = tf.MergeTags(ctx, tx, bannerLink.TagIds, *bannerLink.FeatureId, bannerId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.Status(http.StatusOK)
}

// DeleteBanner удаление баннера
func DeleteBanner(c *gin.Context) {
	//проверка доступа
	isAdmin := c.MustGet("is_admin")
	if !isAdmin.(bool) {
		c.Status(http.StatusForbidden)
		return
	}

	ctx := context.Background()
	conn, tx, err := db.ConnectPoolTrx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка подключения к БД (%v)", err.Error())})
		return
	}
	defer conn.Release()
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)

		} else {
			_ = tx.Commit(ctx)
		}
	}()
	bannerId, err := toInt(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("некорректный идентификатор баннера: %v", err.Error())})
		return
	}
	//проверка существования bannerId
	if err = bn.Exist(context.Background(), conn, bannerId); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	//удаление тегов
	if err = tf.DeleteByBannerId(ctx, tx, bannerId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//удаление контента
	if err = bch.DeleteByBannerId(ctx, tx, bannerId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//удаление баннера
	if err = bn.Delete(ctx, tx, bannerId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func GetToken(c *gin.Context) {
	ctx := context.Background()
	userName, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	user, err := usr.Get(ctx, userName)
	if err != nil || !user.VerifyPassword(password) {
		c.Status(http.StatusUnauthorized)
		return
	} else {
		if token, err := auth.GetToken(user.UserName, user.Password, user.IsAdmin); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		} else {
			c.Writer.Header().Set("token", token)
			c.Status(http.StatusOK)
			return
		}
	}
}

func Auth(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("is_admin", claims["is_admin"])
}

func AuthMiddleware() gin.HandlerFunc {
	return Auth
}
