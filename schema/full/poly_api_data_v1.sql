/*Data for the table `api_poly` */

REPLACE  INTO `api_poly`(`id`,`owner`,`owner_name`,`namespace`,`name`,`title`,`desc`,`access`,`active`,`valid`,`method`,`arrange`,`doc`,`script`,`create_at`,`update_at`,`build_at`,`delete_at`) VALUES 
('poly_AAxfxhEZG2epwoUU1cmAYIKZGyq_xFs3llI4eJcXKMVG','system','系统','/system/poly','permissionInit.p','应用初始化','',0,1,1,'POST','{}','{}','// polyTmpScript_/system/poly/permissionInit.p_2022-03-24T16:43:24CST\nvar _tmp = function(){\n  var d = { \"__input\": __input, } // qyAllLocalData\n\n  d.start = __input.body\n\n  d.start.header = d.start.header || {}\n  d.start._ = [\n    pdCreateNS(\'/system/app\',d.start.appID,\'应用\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID,\'poly\',\'API编排\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID,\'raw\',\'原生API\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID+\'/raw\',\'faas\',\'函数服务\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID+\'/raw\',\'customer\',\'代理第三方API\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID+\'/raw/customer\',\'default\',\'默认分组\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID+\'/raw\',\'inner\',\'平台API\'),\n    pdCreateNS(\'/system/app/\'+d.start.appID+\'/raw/inner\',\'form\',\'表单模型API\'),\n  ]\n  if (true) { // req1, create\n    var _apiPath = format(\"http://structor/api/v1/structor/%v/base/permission/perGroup/create\" ,d.start.appID)\n    var _t = {\n        \"name\": d.start.name,\n        \"description\": d.start.description,\n        \"types\": d.start.types,\n      }\n    var _th = pdNewHttpHeader()\n    pdAddHttpHeader(_th, \"Content-Type\", \"application/json\")\n\n    var _tk = \'\';\n    var _tb = pdAppendAuth(_tk, \'none\', _th, pdToJson(_t))\n    d.req1 = pdToJsobj(\"json\", pdHttpRequest(_apiPath, \"POST\", _tb, _th, pdQueryUser(true)))\n  }\n  d.cond1 = { y: false, }\n  if (d.req1.code==0) {\n    d.cond1.y = true\n    if (true) { // req2, update\n      var _apiPath = format(\"http://structor/api/v1/structor/%v/base/permission/perGroup/update\" ,d.start.appID)\n      var _t = {\n          \"id\": d.req1.data.id,\n          \"scopes\": d.start.scopes,\n        }\n      var _th = pdNewHttpHeader()\n      pdAddHttpHeader(_th, \"Content-Type\", \"application/json\")\n\n      var _tk = \'\';\n      var _tb = pdAppendAuth(_tk, \'none\', _th, pdToJson(_t))\n      d.req2 = pdToJsobj(\"json\", pdHttpRequest(_apiPath, \"POST\", _tb, _th, pdQueryUser(true)))\n    }\n  }\n\n  d.end = {\n    \"createNamespaces\": d.start._,\n    \"req1\": d.req1,\n    \"req2\": sel(d.cond1.y,d.req2,undefined),\n  }\n  return pdToJsonP(d.end)\n}; _tmp();\n',1648111401543,1648111406872,1648111408942,NULL);

/*Data for the table `api_raw` */

REPLACE  INTO `api_raw`(`id`,`owner`,`owner_name`,`namespace`,`name`,`service`,`title`,`desc`,`version`,`path`,`url`,`action`,`method`,`content`,`doc`,`access`,`active`,`valid`,`schema`,`host`,`auth_type`,`create_at`,`update_at`,`delete_at`) VALUES 
('raw_AAAxYvGEh8iIgpjBqleBRjS2J_XKkYJ9IeXyGU9xAt0R','system','系统','/system/form','base_pergroup_create.r','','创建用户组','','last','/api/v1/structor/:appID/base/permission/perGroup/create','http://structor/api/v1/structor/:appID/base/permission/perGroup/create','','POST','{\"x-id\":\"raw_AAAxYvGEh8iIgpjBqleBRjS2J_XKkYJ9IeXyGU9xAt0R\",\"x-action\":\"\",\"x-consts\":[],\"x-input\":{},\"x-output\":{\"body\":{\"type\":\"\",\"name\":\"\",\"data\":null}},\"basePath\":\"/\",\"path\":\"/api/v1/structor/:appID/base/permission/perGroup/create\",\"method\":\"POST\",\"encoding-in\":\"json\",\"encoding-out\":\"json\",\"summary\":\"创建用户组\",\"desc\":\"\"}','{\"x-id\":\"\",\"version\":\"v0.7.3(2021-12-29@f6d9b2b)\",\"x-fmt-inout\":{\"method\":\"POST\",\"url\":\"/api/v1/polyapi/request/system/form/base_pergroup_create.r\",\"input\":{\"inputs\":[{\"type\":\"string\",\"name\":\"X-Polysign-Access-Key-Id\",\"title\":\"签名密钥序号\",\"desc\":\"access_key_id dispatched by poly api server\",\"$appendix$\":true,\"data\":\"KeiIY8098435rty\",\"in\":\"header\",\"mock\":\"KeiIY8098435rty\"},{\"type\":\"string\",\"name\":\"X-Polysign-Timestamp\",\"title\":\"签名时间戳\",\"desc\":\"timestamp format ISO8601: 2006-01-02T15:04:05-0700\",\"$appendix$\":true,\"data\":\"2020-12-31T12:34:56+0800\",\"in\":\"header\",\"mock\":\"2020-12-31T12:34:56+0800\"},{\"type\":\"string\",\"name\":\"X-Polysign-Version\",\"title\":\"签名版本\",\"desc\":\"\\\"1\\\" only current\",\"$appendix$\":true,\"data\":\"1\",\"in\":\"header\",\"mock\":\"1\"},{\"type\":\"string\",\"name\":\"X-Polysign-Method\",\"title\":\"签名方法\",\"desc\":\"\\\"HmacSHA256\\\" only current\",\"$appendix$\":true,\"data\":\"HmacSHA256\",\"in\":\"header\",\"mock\":\"HmacSHA256\"},{\"type\":\"string\",\"name\":\"Access-Token\",\"title\":\"登录授权码\",\"desc\":\"Access-Token from oauth2 if use token access mode\",\"$appendix$\":true,\"data\":null,\"in\":\"header\",\"mock\":\"H3K56789lHIUkjfkslds\"},{\"type\":\"string\",\"name\":\"appID\",\"required\":true,\"data\":null,\"in\":\"path\"},{\"type\":\"object\",\"name\":\"root\",\"data\":[{\"type\":\"string\",\"name\":\"name\",\"data\":null},{\"type\":\"string\",\"name\":\"description\",\"data\":null},{\"type\":\"string\",\"name\":\"x_polyapi_signature\",\"title\":\"参数签名\",\"desc\":\"required if Access-Token doesn\'t use.\\nHmacSHA256 signature of input body: sort query gonic asc|sha256 \\u003cSECRET_KEY\\u003e|base64 std encode\",\"$appendix$\":true,\"data\":\"EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw\"},{\"type\":\"object\",\"name\":\"$polyapi_hide$\",\"title\":\"隐藏参数\",\"desc\":\"polyapi reserved hide args like path args in raw api.\",\"$appendix$\":true,\"data\":[]}],\"in\":\"body\"}]},\"output\":{\"body\":{\"type\":\"\",\"name\":\"\",\"data\":null},\"doc\":[{\"type\":\"object\",\"desc\":\"successful operation\",\"data\":[{\"type\":\"string\",\"name\":\"msg\",\"data\":null},{\"type\":\"number\",\"name\":\"code\",\"data\":null},{\"type\":\"object\",\"name\":\"data\",\"data\":[{\"type\":\"string\",\"name\":\"id\",\"desc\":\"新增后，权限用户组id\",\"data\":null}]}],\"in\":\"body\"}]},\"sampleInput\":[{\"header\":{\"Access-Token\":[\"H3K56789lHIUkjfkslds\"],\"X-Polysign-Access-Key-Id\":[\"KeiIY8098435rty\"],\"X-Polysign-Method\":[\"HmacSHA256\"],\"X-Polysign-Timestamp\":[\"2020-12-31T12:34:56+0800\"],\"X-Polysign-Version\":[\"1\"]},\"body\":{\"$polyapi_hide$\":{\"appID\":\"4x\"},\"description\":\"vYc\",\"name\":\"zxR\",\"x_polyapi_signature\":\"EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw\"}},{\"header\":{\"登录授权码\":[\"H3K56789lHIUkjfkslds\"],\"签名密钥序号\":[\"KeiIY8098435rty\"],\"签名方法\":[\"HmacSHA256\"],\"签名时间戳\":[\"2020-12-31T12:34:56+0800\"],\"签名版本\":[\"1\"]},\"body\":{\"description\":\"u_vOjzgB\",\"name\":\"Kk4L0opORX\",\"参数签名\":\"EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw\",\"隐藏参数\":{\"appID\":\"ovK-GRJ28N\"}}}],\"sampleOutput\":[{\"resp\":{\"code\":17,\"data\":{\"id\":\"t8\"},\"msg\":\"xW4\"}},{\"resp\":{\"code\":4,\"data\":{\"id\":\"8u9sOoVCP\"},\"msg\":\"fYSmb134qiG\"}}]},\"x-swagger\":{\"x-consts\":null,\"host\":\"structor\",\"swagger\":\"2.0\",\"info\":{\"title\":\"\",\"version\":\"last\",\"description\":\"auto generated\",\"contact\":{\"name\":\"\",\"url\":\"\",\"email\":\"\"}},\"schemes\":[\"http\"],\"basePath\":\"/\",\"paths\":{\"/api/v1/structor/:appID/base/permission/perGroup/create\":{\"post\":{\"x-consts\":[],\"operationId\":\"base_pergroup_create\",\"parameters\":[{\"name\":\"appID\",\"in\":\"path\",\"description\":\"\",\"required\":true,\"type\":\"string\"},{\"name\":\"root\",\"in\":\"body\",\"schema\":{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"name\":{\"type\":\"string\"},\"description\":{\"type\":\"string\"}},\"required\":[]}}],\"responses\":{\"200\":{\"description\":\"successful operation\",\"schema\":{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"code\":{\"type\":\"number\"},\"data\":{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"string\",\"description\":\"新增后，权限用户组id\"}}},\"msg\":{\"type\":\"string\"}},\"required\":[\"code\"]}}},\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"],\"summary\":\"创建用户组\",\"description\":\"\"}}}}}',0,1,1,'http','structor','none',1648111398882,1648111398882,NULL),
('raw_AM51O-2rUXb1RVDXnOvAo8FRq5BJzaO4vdO8QVx-qZ1n','system','系统','/system/form','base_pergroup_update.r','','给用户组加入人员或者部门','','last','/api/v1/structor/:appID/base/permission/perGroup/update','http://structor/api/v1/structor/:appID/base/permission/perGroup/update','','POST','{\"x-id\":\"raw_AM51O-2rUXb1RVDXnOvAo8FRq5BJzaO4vdO8QVx-qZ1n\",\"x-action\":\"\",\"x-consts\":[],\"x-input\":{},\"x-output\":{\"body\":{\"type\":\"\",\"name\":\"\",\"data\":null}},\"basePath\":\"/\",\"path\":\"/api/v1/structor/:appID/base/permission/perGroup/update\",\"method\":\"POST\",\"encoding-in\":\"json\",\"encoding-out\":\"json\",\"summary\":\"给用户组加入人员或者部门\",\"desc\":\"\"}','{\"x-id\":\"\",\"version\":\"v0.7.3(2021-12-29@f6d9b2b)\",\"x-fmt-inout\":{\"method\":\"POST\",\"url\":\"/api/v1/polyapi/request/system/form/base_pergroup_update.r\",\"input\":{\"inputs\":[{\"type\":\"string\",\"name\":\"X-Polysign-Access-Key-Id\",\"title\":\"签名密钥序号\",\"desc\":\"access_key_id dispatched by poly api server\",\"$appendix$\":true,\"data\":\"KeiIY8098435rty\",\"in\":\"header\",\"mock\":\"KeiIY8098435rty\"},{\"type\":\"string\",\"name\":\"X-Polysign-Timestamp\",\"title\":\"签名时间戳\",\"desc\":\"timestamp format ISO8601: 2006-01-02T15:04:05-0700\",\"$appendix$\":true,\"data\":\"2020-12-31T12:34:56+0800\",\"in\":\"header\",\"mock\":\"2020-12-31T12:34:56+0800\"},{\"type\":\"string\",\"name\":\"X-Polysign-Version\",\"title\":\"签名版本\",\"desc\":\"\\\"1\\\" only current\",\"$appendix$\":true,\"data\":\"1\",\"in\":\"header\",\"mock\":\"1\"},{\"type\":\"string\",\"name\":\"X-Polysign-Method\",\"title\":\"签名方法\",\"desc\":\"\\\"HmacSHA256\\\" only current\",\"$appendix$\":true,\"data\":\"HmacSHA256\",\"in\":\"header\",\"mock\":\"HmacSHA256\"},{\"type\":\"string\",\"name\":\"Access-Token\",\"title\":\"登录授权码\",\"desc\":\"Access-Token from oauth2 if use token access mode\",\"$appendix$\":true,\"data\":null,\"in\":\"header\",\"mock\":\"H3K56789lHIUkjfkslds\"},{\"type\":\"string\",\"name\":\"appID\",\"required\":true,\"data\":null,\"in\":\"path\"},{\"type\":\"object\",\"name\":\"root\",\"title\":\"empty object\",\"data\":[{\"type\":\"string\",\"name\":\"id\",\"desc\":\"用户组权限id\",\"data\":null},{\"type\":\"array\",\"name\":\"scopes\",\"data\":[{\"type\":\"object\",\"name\":\"\",\"data\":[{\"type\":\"number\",\"name\":\"type\",\"desc\":\"1 人员 2 部门\",\"data\":null},{\"type\":\"string\",\"name\":\"id\",\"desc\":\"人员或者部门id\",\"data\":null},{\"type\":\"string\",\"name\":\"name\",\"desc\":\"人员或者部门名字\",\"data\":null}]}]},{\"type\":\"string\",\"name\":\"x_polyapi_signature\",\"title\":\"参数签名\",\"desc\":\"required if Access-Token doesn\'t use.\\nHmacSHA256 signature of input body: sort query gonic asc|sha256 \\u003cSECRET_KEY\\u003e|base64 std encode\",\"$appendix$\":true,\"data\":\"EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw\"},{\"type\":\"object\",\"name\":\"$polyapi_hide$\",\"title\":\"隐藏参数\",\"desc\":\"polyapi reserved hide args like path args in raw api.\",\"$appendix$\":true,\"data\":[]}],\"in\":\"body\"}]},\"output\":{\"body\":{\"type\":\"\",\"name\":\"\",\"data\":null},\"doc\":[{\"type\":\"object\",\"desc\":\"successful operation\",\"data\":[{\"type\":\"object\",\"name\":\"data\",\"data\":[]},{\"type\":\"string\",\"name\":\"msg\",\"data\":null},{\"type\":\"number\",\"name\":\"code\",\"data\":null}],\"in\":\"body\"}]},\"sampleInput\":[{\"header\":{\"Access-Token\":[\"H3K56789lHIUkjfkslds\"],\"X-Polysign-Access-Key-Id\":[\"KeiIY8098435rty\"],\"X-Polysign-Method\":[\"HmacSHA256\"],\"X-Polysign-Timestamp\":[\"2020-12-31T12:34:56+0800\"],\"X-Polysign-Version\":[\"1\"]},\"body\":{\"$polyapi_hide$\":{\"appID\":\"z4YpHB\"},\"id\":\"Ing5IkP2\",\"scopes\":[{\"id\":\"cEDmW\",\"name\":\"vd4cGe2JAl\",\"type\":15}],\"x_polyapi_signature\":\"EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw\"}},{\"header\":{\"登录授权码\":[\"H3K56789lHIUkjfkslds\"],\"签名密钥序号\":[\"KeiIY8098435rty\"],\"签名方法\":[\"HmacSHA256\"],\"签名时间戳\":[\"2020-12-31T12:34:56+0800\"],\"签名版本\":[\"1\"]},\"body\":{\"id\":\"x010u\",\"scopes\":[{\"id\":\"Kp3\",\"name\":\"3tBMriHYO\",\"type\":16}],\"参数签名\":\"EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw\",\"隐藏参数\":{\"appID\":\"pfi\"}}}],\"sampleOutput\":[{\"resp\":{\"code\":10,\"data\":{},\"msg\":\"aEK\"}},{\"resp\":{\"code\":9,\"data\":{},\"msg\":\"ZlHhTVQRIE\"}}]},\"x-swagger\":{\"x-consts\":null,\"host\":\"structor\",\"swagger\":\"2.0\",\"info\":{\"title\":\"\",\"version\":\"last\",\"description\":\"auto generated\",\"contact\":{\"name\":\"\",\"url\":\"\",\"email\":\"\"}},\"schemes\":[\"http\"],\"basePath\":\"/\",\"paths\":{\"/api/v1/structor/:appID/base/permission/perGroup/update\":{\"post\":{\"x-consts\":[],\"operationId\":\"base_pergroup_update\",\"parameters\":[{\"name\":\"appID\",\"in\":\"path\",\"description\":\"\",\"required\":true,\"type\":\"string\"},{\"name\":\"root\",\"in\":\"body\",\"schema\":{\"type\":\"object\",\"title\":\"empty object\",\"properties\":{\"id\":{\"type\":\"string\",\"description\":\"用户组权限id\"},\"scopes\":{\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"type\":{\"type\":\"integer\",\"description\":\"1 人员 2 部门\"},\"id\":{\"type\":\"string\",\"description\":\"人员或者部门id\"},\"name\":{\"type\":\"string\",\"description\":\"人员或者部门名字\"}},\"required\":[\"type\",\"id\",\"name\"]}}},\"required\":[\"id\",\"scopes\"]}}],\"responses\":{\"200\":{\"description\":\"successful operation\",\"schema\":{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"code\":{\"type\":\"number\"},\"data\":{\"type\":\"object\",\"properties\":{}},\"msg\":{\"type\":\"string\"}}}}},\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"],\"summary\":\"给用户组加入人员或者部门\",\"description\":\"\"}}}}}',0,1,1,'http','structor','none',1648111400071,1648111400071,NULL);
