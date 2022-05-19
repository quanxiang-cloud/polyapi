// json => js obj
function pdFromJson(jsonStr){
    return JSON.parse(jsonStr)
}

// js obj => json
function pdToJson(jsObj){
    return JSON.stringify(jsObj)
}

// js obj => json, pretty style
function pdToJsonP(jsObj){
    return JSON.stringify(jsObj, undefined, 2)
}

// two way select logic
function pdSelect(cond, yes, no){
	if (cond) {
		return yes
	}else{
		return no
	}
}

function sel(cond, yes, no){
	return pdSelect(cond, yes, no)
}

// filt filed for object
function pdFiltObject(input, config, filtFunc){
	if (input instanceof Array){
		return _filtArray(input, config, filtFunc);
	}else{
		return _filtObject(input, config);
	}
}

// check if an object is empty
function _isEmptyObject(obj) {
　　for (var k in obj){
　　　　return false;
　　}　　
　　return true;
}　

// filt filed for object
function _filtObject(input, config){
	if (typeof input !== 'object'){
		return input;
	}
	
	var ret = {};
	var w = config.white || {}
	var b = config.black || {}
	var bIsEmpty = _isEmptyObject(b)
	for (var k in input) {
		var newName = w[k];
		var deny = b[k];
		if ((!bIsEmpty && deny === undefined) || newName !== undefined){
			if ((typeof newName !== 'string') || newName === ""){
				newName = k;
			}
			ret[newName] = input[k];
		}
	}
	return ret;
}

// filt filed and data for array
function _filtArray(input, config, filtFunc){
	var ret = [];
	for (var i = 0; i < input.length; i++) {
		var elem = input[i];
		if (filtFunc===undefined || filtFunc(i,elem)){
			var x = _filtObject(elem, config);
			ret.push(x);
		}
	}
	return ret;
}

// pdMergeObjs(jsObj...) obj
function pdMergeObjs(){
	var ret = {};
	var i = 0;
    for(var i; i<arguments.length; i++){
       var obj = arguments[i];
       if (typeof obj === 'object'){
       	for (var k in obj) {
       		ret[k] = obj[k];
       	}
       }
    }
    return ret;
}

// pdToJsobj("encoding", strData) obj
function pdToJsobj(encoding, strData){
	switch (encoding) {
	case "json":
		return pdFromJson(strData)
		break
	case "xml":
		var obj = pdFromXml(strData)
		var jsn = pdToJson(obj)
		return pdFromJson(jsn)
		break
	case "yaml":
		var obj = pdFromYaml(strData)
		var jsn = pdToJson(obj)
		return pdFromJson(jsn)
		break
	default:
		break
	}
	return {}
}
