package esphomeproto_test

import (
	"testing"

	"github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/matryer/is"
)

func TestCoreEntityCategoryEnum(t *testing.T) {
	is := is.New(t)
	is.Equal(int(entity.CategoryNone), int(esphomeproto.EntityCategory_ENTITY_CATEGORY_NONE))
	is.Equal(int(entity.CategoryConfig), int(esphomeproto.EntityCategory_ENTITY_CATEGORY_CONFIG))
	is.Equal(int(entity.CategoryDisgnostic), int(esphomeproto.EntityCategory_ENTITY_CATEGORY_DIAGNOSTIC))
	is.Equal(len(esphomeproto.EntityCategory_name), 3)
}
