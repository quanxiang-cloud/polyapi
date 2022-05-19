package xml

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestXml(t *testing.T) {
	enc := Encoder{}
	var xx = []interface{}{
		123,
		"string",
		true,
		12.34,
		[]string{"a", "b", "c"},
		[]int{1, 2, 3, 4},
		map[string]interface{}{"a": 1, "b": "xx", "c": true, "d": []string{"x", "y"}},
		struct {
			X int
			Y string
			Z []string
		}{1, "yy", []string{"x1", "y1"}},
	}
	b, err := enc.Encode(xx, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("xml:\n%s\n", string(b))

	dec := Decoder{}
	d, err := dec.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("decode xml:\n%#v\n", d)

	tt := []string{`
	<?xml version="1.0" encoding="UTF-8"?>
<data>

 <FRegionId_0>41680</FRegionId_0> 
 <FRegionId_1>41680</FRegionId_1> 
 <FUsedCount_0>0</FUsedCount_0> 
 <FUsedCount_1>0</FUsedCount_1> 
 <Faddrstreet_0>山东潍坊</Faddrstreet_0> 
 <Faddrstreet_1>腾讯大厦</Faddrstreet_1> 
 <Fcreate_time_0>2011-07-08 19:37:08</Fcreate_time_0> 
 <Fcreate_time_1>2011-07-08 19:36:23</Fcreate_time_1> 
 <Findex_0>2</Findex_0> 
 <Findex_1>1</Findex_1> 
 <Flastuse_time_0>2011-07-08 19:37:08</Flastuse_time_0> 
 <Flastuse_time_1>2011-07-08 19:36:23</Flastuse_time_1> 
 <Fmobile_0>15812345678</Fmobile_0> 
 <Fmobile_1>18612345678</Fmobile_1> 
 <Fmod_time_0>2011-07-08 19:37:08</Fmod_time_0> 
 <Fmod_time_1>2011-07-08 19:36:23</Fmod_time_1> 
 <Fname_0>张三</Fname_0> 
 <Fname_1>张三</Fname_1> 	
 <Ftel_0>05361234567</Ftel_0>  	
 <Ftel_1>07551234567</Ftel_1>  
 <Fzipcode_0>276323</Fzipcode_0> 
 <Fzipcode_1>510640</Fzipcode_1> 
 <msg>ok</msg> 
 <ret>0</ret> 
 <ret_num>2</ret_num> 
</data>
	`,
		`
	<?xml version="1.0" encoding="UTF-8"?>
<data>
  <ret>0</ret>
  <errcode>0</errcode>
  <msg>ok</msg>
  <data>
     <timestamp>128679200</timestamp>
     <hasnext>0</hasnext>
     <totalnum>2</totalnum>
     <info>
          <text></text>
          <origtext></origtext>
          <count>2</count>
          <from>来自网页</from>
          <id>7987543214334</id>
          <image></image>
          <name>abc</name>
          <openid>B624064BA065E01CB73F835017FE96FA</openid>
          <nick>abcd</nick>
          <self>0</self>
          <timestamp>1285813236</timestamp>
          <type>1</type>
          <head>http://app.qlogo.cn/mbloghead/563ad8b6be488a07a694</head>
          <location>广东 深圳</location>
          <country_code>1</country_code>
          <province_code>44</province_code>		   
          <city_code>3</city_code>
          <isvip>0</isvip>
          <geo>null</geo>		   
     </info>
  </data>
  <user>
     <name:nick></name:nick>
  </user>   
</data>
	`,
	}
	for _, v := range tt {
		dec := Decoder{}
		d, err := dec.Decode([]byte(v))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("decode xml:\n%#v\n", d)
		b, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
