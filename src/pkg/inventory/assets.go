package inventory

import (
	"fmt"
	"reflect"
)

type AssetType string
type IdProperty string

const (
	AssetTypeOrganization AssetType = "organization"
	AssetTypeProject      AssetType = "project"
	AssetTypeArtefact     AssetType = "artefact"
	AssetTypeRepository   AssetType = "repository"
	AssetTypeFqdn         AssetType = "fqdn"
	AssetTypeIp           AssetType = "ip"
	AssetTypeEndpoint     AssetType = "endpoint"
	AssetTypeDomain       AssetType = "domain"
	AssetTypeHostname     AssetType = "hostname"
	AssetAccount          AssetType = "account"
	AssetPackage          AssetType = "package"
	AssetTypeHash         AssetType = "hash"
	AssetTypeCertificate  AssetType = "certificate"
	AssetTypeSecret       AssetType = "secret"
)

var (
	AllAssetTypes = []AssetType{
		AssetTypeOrganization,
		AssetTypeProject,
		AssetTypeArtefact,
		AssetTypeRepository,
		AssetTypeFqdn,
		AssetTypeIp,
		AssetTypeEndpoint,
		AssetTypeDomain,
		AssetTypeHostname,
		AssetAccount,
		AssetPackage,
		AssetTypeHash,
		AssetTypeCertificate,
		AssetTypeSecret,
	}
	// empty asset parents is equals to no parent
	// left is parent right is their childrens
	ParentMap = map[AssetType][]AssetType{
		AssetTypeOrganization: {},
		AssetTypeProject:      {AssetTypeOrganization},
		AssetTypeRepository:   {AssetTypeProject},
		AssetTypeArtefact:     {AssetTypeRepository},
		AssetPackage:          {AssetTypeArtefact, AssetTypeRepository},
		AssetTypeHash:         {AssetTypeArtefact, AssetPackage},
		AssetTypeFqdn:         {AssetTypeDomain},
		AssetTypeIp:           {AssetTypeFqdn},
		AssetTypeEndpoint:     {AssetTypeFqdn},
		AssetTypeDomain:       {AssetTypeOrganization},
		AssetTypeHostname:     {AssetTypeIp},
		AssetAccount:          {AssetTypeOrganization},
		AssetTypeCertificate:  {AssetTypeOrganization},
		AssetTypeSecret:       {AssetTypeArtefact, AssetPackage, AssetTypeRepository},
	}
)

type Asset[
	T AssetType,
	TId IdProperty,
	TSource interface{},
] struct {
	Kind       T
	IdProperty TId
	Origin     string
	Source     TSource
}

// type Asset struct {
// 	Kind       AssetType   `json:"kind,omitempty"`
// 	IdProperty string      `json:"idProperty"`
// 	Origin     string      `json:"origin"`
// 	Source     interface{} `json:"source"`
// }

func (a Asset[T, TId, TSource]) Type() T {
	return a.Kind
}
func (a Asset[T, TId, TSource]) String() string {
	return fmt.Sprintf("%s:%s", a.Kind, a.Id())
}
func (a *Asset[T, TId, TSource]) IsValid() bool {
	var valid bool = false
	for _, t := range AllAssetTypes {
		if t == AssetType(a.Kind) {
			valid = true
			break
		}
	}
	valid = valid && a.IdProperty != "" && !reflect.ValueOf(a.Source).IsNil() && a.Id() != "" && a.Origin != ""
	return valid
}
func (a *Asset[T, TId, TSource]) Signature() string {
	return Checksum(fmt.Sprintf("%s.%s.%s", a.Kind, a.Origin, a.Id()))
}
func (a *Asset[T, TId, TSource]) Id() string {
	// id := a.Source.(map[string]interface{})[a.IdProperty]
	id := reflect.ValueOf(a.Source).FieldByName(string(a.IdProperty)).Interface()	
	if id == nil {
		return ""
	}
	return id.(string)
}
func NewAsset[
	T AssetType,
	TId IdProperty,
	TSource interface{},
](
	kind T,
	idProperty TId,
	origin string,
	source TSource,
) Asset[T, TId, TSource] {
	return Asset[T, TId, TSource]{
		Kind:       kind,
		IdProperty: idProperty,
		Origin:     origin,
		Source:     source,
	}
}
type OrganizationAssetSource struct {
	Id   string `json:"id,omitempty"`
}
type OrganizationAsset struct {
	Asset[AssetType, IdProperty, OrganizationAssetSource] `json:",inline,omitempty"`
}
func NewOrganizationAsset(args struct {
	Name string
}) *OrganizationAsset {
	return &OrganizationAsset{
		Asset: NewAsset(AssetTypeOrganization, "id", "internal", OrganizationAssetSource{Id: "1"}),
	}
}
// type OrganizationAsset struct {
// 	Asset `json:",inline,omitempty"`
// 	Name  string `json:"name,omitempty"`
// }

// func NewOrganizationAsset(args struct {
// 	Name string
// }) *OrganizationAsset {
// 	asset := NewAsset(AssetTypeOrganization, "id", "internal", map[string]interface{}{})
// 	return &OrganizationAsset{
// 		Asset: asset,
// 		Name:  args.Name,
// 	}
// }

// type ProjectAsset struct {
// 	Asset `json:",inline,omitempty"`
// 	Name  string `json:"name,omitempty"`
// }
