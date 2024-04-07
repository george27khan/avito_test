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
	TagIds    []int64 `json:"tag_ids"`
	FeatureId int64   `json:"feature_id"`
	Content   string  `json:"content"`
	IsActive  bool    `json:"is_active"`
}

// PostBanner добавление баннера
func PostBanner(c *gin.Context) {
	var body PostBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	banner := bn.Banner{Content: []byte(body.Content),
		IsActive:  body.IsActive,
		FeatureId: body.FeatureId,
		Tags:      body.TagIds,
	}

	if err := banner.Check(context.Background()); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	idBanner, err := banner.Insert(context.Background())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusCreated, gin.H{"banner_id": idBanner, "answer": http.StatusInternalServerError})
}

func toInt64(val string) (int64, error) {
	res, err := strconv.ParseInt(val, 10, 64)
	return res, err
}

// GetBanner добавление баннера
func GetBanner(c *gin.Context) {
	var (
		body       PostBody
		bannerList []bn.Banner
	)
	feature := c.Query("feature_id")
	tag := c.Query("tag_id")
	limit := c.Query("limit")
	offset := c.Query("offset")
	token := c.GetHeader("token")

	fmt.Println(feature, tag, limit, offset, token)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if tag != "" && feature != "" {
		tagInt, err1 := toInt64(tag)
		featureInt, err2 := toInt64(feature)
		if err1 == nil && err2 == nil {
			if banner, err := bn.Get(context.Background(), tagInt, featureInt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				bannerList = append(bannerList, banner)
			}
		}
	}

	//if err := banner.Ceck(context.Background()); err != nil {
	//	c.JSON(http.StatusBadRequest, err.Error())
	//}
	//idBanner, err := banner.Insert(context.Background())
	//if err != nil {
	//	c.String(http.StatusInternalServerError, err.Error())
	//}
	//c.JSON(http.StatusCreated, gin.H{"banner_id": idBanner, "answer": http.StatusInternalServerError})

}
