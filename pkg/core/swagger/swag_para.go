package swagger

// SwagParameters represents input from swagger
type SwagParameters []SwagValue

// SwagValue represents common value structure in swagger
type SwagValue map[string]interface{}

// SwagObjectPropperties represents propperties of object in swagger
type SwagObjectPropperties map[string]SwagValue
