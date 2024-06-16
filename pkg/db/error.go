package db

import "errors"

// Обозначение ошибок
var (
	ErrMigrate          = errors.New("migration failed")
	ErrDuplicate        = errors.New("record already exists")
	ErrNotExist         = errors.New("row does not exist")
	ErrUpdateFailed     = errors.New("update failed")
	ErrDeleteFailed     = errors.New("delete failed")
	ErrParentRegion     = errors.New("some region has this parent region")
	ErrRegionType       = errors.New("some region has this region type")
	ErrParamNotFound    = errors.New("param not found")
	ErrAuthorize        = errors.New("authorize failed")
	ErrValidate         = errors.New("validate failed")
	ErrExistCoordinates = errors.New("this coordinates is already exist")
)
