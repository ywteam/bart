package inventory

// import (
// 	"testing"
// )

// func TestAsset_IsValid(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		asset  Asset
// 		expect bool
// 	}{
// 		{
// 			name: "valid asset",
// 			asset: Asset{
// 				Kind:       AssetTypeOrganization,
// 				IdProperty: "id",
// 				Origin:     "internal",
// 				Source: map[string]interface{}{
// 					"id": "1",
// 				},
// 			},
// 			expect: true,
// 		},
// 		{
// 			name: "invalid asset type",
// 			asset: Asset{
// 				Kind:       "invalid",
// 				IdProperty: "id",
// 				Origin:     "internal",
// 				Source: map[string]interface{}{
// 					"id": "1",
// 				},
// 			},
// 			expect: false,
// 		},
// 		{
// 			name: "missing id property",
// 			asset: Asset{
// 				Kind:       AssetTypeOrganization,
// 				IdProperty: "",
// 				Origin:     "internal",
// 				Source: map[string]interface{}{
// 					"id": "1",
// 				},
// 			},
// 			expect: false,
// 		},
// 		{
// 			name: "missing source",
// 			asset: Asset{
// 				Kind:       AssetTypeOrganization,
// 				IdProperty: "id",
// 				Origin:     "internal",
// 				Source:     nil,
// 			},
// 			expect: false,
// 		},
// 		{
// 			name: "missing id in source",
// 			asset: Asset{
// 				Kind:       AssetTypeOrganization,
// 				IdProperty: "id",
// 				Origin:     "internal",
// 				Source: map[string]interface{}{
// 					"name": "My Organization",
// 				},
// 			},
// 			expect: false,
// 		},
// 		{
// 			name: "missing origin",
// 			asset: Asset{
// 				Kind:       AssetTypeOrganization,
// 				IdProperty: "id",
// 				Origin:     "",
// 				Source: map[string]interface{}{
// 					"id": "1",
// 				},
// 			},
// 			expect: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.asset.IsValid(); got != tt.expect {
// 				t.Errorf("Asset.IsValid() = %v, expect %v", got, tt.expect)
// 			}
// 		})
// 	}
// }

// func TestAsset_Signature(t *testing.T) {
// 	asset := Asset{
// 		Kind:       AssetTypeOrganization,
// 		IdProperty: "id",
// 		Origin:     "internal",
// 		Source: map[string]interface{}{
// 			"id": "1",
// 		},
// 	}
// 	expectedSignature := Checksum("organization.internal.1")
// 	if got := asset.Signature(); got != expectedSignature {
// 		t.Errorf("Asset.Signature() = %v, expect %v", got, expectedSignature)
// 	}
// }

// func TestAsset_Id(t *testing.T) {
// 	asset := Asset{
// 		Kind:       AssetTypeOrganization,
// 		IdProperty: "id",
// 		Origin:     "internal",
// 		Source: map[string]interface{}{
// 			"id": "1",
// 		},
// 	}
// 	expectedId := "1"
// 	if got := asset.Id(); got != expectedId {
// 		t.Errorf("Asset.Id() = %v, expect %v", got, expectedId)
// 	}
// }

// func TestAsset_OrganizationAsset(t *testing.T) {
// 	org := OrganizationAsset{
// 		Asset: Asset{
// 			Kind:       AssetTypeOrganization,
// 			Origin:     "internal",
// 			IdProperty: "id",
// 			Source: map[string]interface{}{
// 				"id": "1",
// 			},
// 		},
// 		Name: "My Organization",
// 	}
// 	expectedSignature := Checksum("organization.internal.1")
// 	if got := org.Signature(); got != expectedSignature {
// 		t.Errorf("OrganizationAsset.Signature() = %v, expect %v", got, expectedSignature)
// 	}
// 	t.Log("OrganizationAsset: ", JsonStringfy(org, true))
// }
// func TestAsset_ProjectAsset(t *testing.T) {
// 	project := ProjectAsset{
// 		Asset: Asset{
// 			Kind:       AssetTypeProject,
// 			Origin:     "internal",
// 			IdProperty: "id",
// 			Source: map[string]interface{}{
// 				"id": "2",
// 			},
// 		},
// 		Name: "My Project",
// 	}
// 	expectedSignature := Checksum("project.internal.2")
// 	if got := project.Signature(); got != expectedSignature {
// 		t.Errorf("ProjectAsset.Signature() = %v, expect %v", got, expectedSignature)
// 	}
// 	t.Log("ProjectAsset: ", JsonStringfy(project, true))
// }
