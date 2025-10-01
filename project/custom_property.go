package project

func InsertOrUpdateEnumType(proj *TiledProject, definitions ...TiledEnumPropertyType) error {
	nextID := getNextPropertyTypeID(proj)

	for _, def := range definitions {
		existing := getExistingEnumType(proj, def.Name)

		if existing != nil {
			def.ID = existing.ID
			*existing = def
			continue
		}

		def.ID = nextID
		nextID++

		proj.EnumPropertyTypes = append(proj.EnumPropertyTypes, def)
	}

	return nil
}

func InsertOrUpdateClassType(proj *TiledProject, definitions ...TiledClassPropertyType) error {
	nextID := getNextPropertyTypeID(proj)

	for _, def := range definitions {
		existing := getExistingClassType(proj, def.Name)

		if existing != nil {
			def.ID = existing.ID
			*existing = def
			continue
		}

		def.ID = nextID
		nextID++

		proj.ClassPropertyTypes = append(proj.ClassPropertyTypes, def)
	}

	return nil
}

func getNextPropertyTypeID(proj *TiledProject) int {
	maxID := 0
	for _, enumType := range proj.EnumPropertyTypes {
		if enumType.ID > maxID {
			maxID = enumType.ID
		}
	}
	for _, classType := range proj.ClassPropertyTypes {
		if classType.ID > maxID {
			maxID = classType.ID
		}
	}
	return maxID + 1
}

func getExistingEnumType(proj *TiledProject, name string) *TiledEnumPropertyType {
	for i, enumType := range proj.EnumPropertyTypes {
		if enumType.Name == name {
			return &proj.EnumPropertyTypes[i]
		}
	}
	return nil
}

func getExistingClassType(proj *TiledProject, name string) *TiledClassPropertyType {
	for i, classType := range proj.ClassPropertyTypes {
		if classType.Name == name {
			return &proj.ClassPropertyTypes[i]
		}
	}
	return nil
}
