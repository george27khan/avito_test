package banner_handler

import (
	bn "avito_test/pkg/db_avito_banner/banner"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostBody struct {
	TagIds    []int64     `json:"tag_ids"`
	FeatureId int64       `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
}

func toInt(val string) (int64, error) {
	res, err := strconv.ParseInt(val, 10, 64)
	return res, err
}

// PostBanner добавление баннера
func PostBanner(c *gin.Context) {
	var banner bn.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := banner.Check(context.Background()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bannerId, err := banner.Create(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"banner_id": bannerId})
	return
}

// GetBanner добавление баннера
func GetBanner(c *gin.Context) {
	var (
		bannerList []bn.Banner
	)

	feature := c.Query("feature_id")
	tag := c.Query("tag_id")
	limit := c.Query("limit")
	offset := c.Query("offset")
	//token := c.GetHeader("token")

	if tag != "" && feature != "" {
		tagInt, err1 := toInt(tag)
		featureInt, err2 := toInt(feature)
		if err1 == nil && err2 == nil {
			if banner, err := bn.Get(context.Background(), tagInt, featureInt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				bannerList = append(bannerList, banner)
				c.JSON(http.StatusOK, bannerList)
				return
			}
		}
	} else if tag != "" {
		var (
			tagInt, limitInt, offsetInt int64
			err                         error
		)
		if tagInt, err = toInt(tag); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный tag_id=%v - %v", tag, err.Error())})
			return
		}
		if limitInt, err = toInt(limit); err != nil && limit != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный limit=%v - %v", limit, err.Error())})
			return
		}
		if offsetInt, err = toInt(offset); err != nil && offset != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный offset=%v - %v", offset, err.Error())})
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
		var (
			featureInt, limitInt, offsetInt int64
			err                             error
		)
		if featureInt, err = toInt(feature); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный feature_id=%v - %v", feature, err.Error())})
			return
		}
		if limitInt, err = toInt(limit); err != nil && limit != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный limit=%v - %v", limit, err.Error())})
			return
		}
		if offsetInt, err = toInt(offset); err != nil && offset != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный offset=%v - %v", offset, err.Error())})
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

// PatchBanner изменение баннера
func PatchBanner(c *gin.Context) {
	var banner bn.Banner
	bannerId, err := toInt(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//проверка существования bannerId
	if err := bn.Exist(context.Background(), bannerId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//обновление данных баннера
	if err := banner.Update(context.Background(), bannerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//token := c.GetHeader("token")

	//if tag != "" && feature != "" {
	//	tagInt, err1 := toInt(tag)
	//	featureInt, err2 := toInt(feature)
	//	if err1 == nil && err2 == nil {
	//		if banner, err := bn.Get(context.Background(), tagInt, featureInt); err != nil {
	//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//			return
	//		} else {
	//			bannerList = append(bannerList, banner)
	//			c.JSON(http.StatusOK, bannerList)
	//			return
	//		}
	//	}
	//} else if tag != "" {
	//	var (
	//		tagInt, limitInt, offsetInt int64
	//		err                         error
	//	)
	//	if tagInt, err = toInt(tag); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный tag_id=%v - %v", tag, err.Error())})
	//		return
	//	}
	//	if limitInt, err = toInt(limit); err != nil && limit != "" {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный limit=%v - %v", limit, err.Error())})
	//		return
	//	}
	//	if offsetInt, err = toInt(offset); err != nil && offset != "" {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный offset=%v - %v", offset, err.Error())})
	//		return
	//	}
	//	if bannerList, err := bn.GetByTag(context.Background(), tagInt, limitInt, offsetInt); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	} else {
	//		c.JSON(http.StatusOK, bannerList)
	//		return
	//	}
	//} else if feature != "" {
	//	var (
	//		featureInt, limitInt, offsetInt int64
	//		err                             error
	//	)
	//	if featureInt, err = toInt(feature); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный feature_id=%v - %v", feature, err.Error())})
	//		return
	//	}
	//	if limitInt, err = toInt(limit); err != nil && limit != "" {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный limit=%v - %v", limit, err.Error())})
	//		return
	//	}
	//	if offsetInt, err = toInt(offset); err != nil && offset != "" {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Передан некорректный offset=%v - %v", offset, err.Error())})
	//		return
	//	}
	//	if bannerList, err := bn.GetByFeature(context.Background(), featureInt, limitInt, offsetInt); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		return
	//	} else {
	//		c.JSON(http.StatusOK, bannerList)
	//		return
	//	}
	//}

}
