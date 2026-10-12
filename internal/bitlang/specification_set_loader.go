package bitlang

//	{
//	  責務: [
//	    LoadSpecificationSet: 意味解析に必要な仕様群を読み込み横断検証する
//	  ]
//	  処理: [
//	    1: プロパティ仕様を厳密に読み込む
//	    2: 宣言仕様を厳密に読み込む
//	    3: 両仕様間の参照整合性を確認する
//	  ]
//	  引数: [
//	    propertiesData: プロパティ仕様の JSON バイト列
//	    declarationsData: 宣言仕様の JSON バイト列
//	  ]
//	  戻り値: [
//	    specs: 個別・横断検証済みの仕様群
//	    error: 読込または整合性検証の失敗
//	  ]
//	}
func LoadSpecificationSet(propertiesData []byte, declarationsData []byte) (SpecificationSet, error) {
	properties, err := LoadPropertySpecification(propertiesData)
	if err != nil {
		return SpecificationSet{}, err
	}
	declarations, err := LoadDeclarationSpecification(declarationsData)
	if err != nil {
		return SpecificationSet{}, err
	}
	if err := ValidateSpecificationConsistency(properties, declarations); err != nil {
		return SpecificationSet{}, err
	}
	return SpecificationSet{
		Properties:   properties,
		Declarations: declarations,
	}, nil
}
