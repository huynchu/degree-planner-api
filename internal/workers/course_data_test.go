package workers

// func (p Prerequisite) TransformPrereqRecursive(res *[][]string) []string {
// 	if p.Type == "and" {
// 		for _, prereq := range p.Nested {
// 			*res = append(*res, prereq.TransformPrereqRecursive(res))
// 		}
// 		return []string{}
// 	} else if p.Type == "or" {
// 		or := []string{}
// 		for _, prereq := range p.Nested {
// 			or = append(or, prereq.TransformPrereqRecursive(res)...)
// 		}
// 		return or
// 	} else {
// 		return []string{p.Course}
// 	}
// }

// func (p Prerequisite) TransformPrereq() [][]string {
// 	if p.Type == "course" {
// 		return [][]string{{p.Course}}
// 	} else {
// 		res := [][]string{}
// 		p.TransformPrereqRecursive(&res)
// 		return res
// 	}
// }

// func TestParseCoursePreReqs(t *testing.T) {
// 	jsonFile, err := os.Open("prereq_data.json")
// 	if err != nil {
// 		t.Errorf("error opening json file: %v", err)
// 	}

// 	defer jsonFile.Close()

// 	dataBytes, _ := io.ReadAll(jsonFile)

// 	coursePrereqDataMap := make(map[string]interface{})
// 	err = json.Unmarshal(dataBytes, &coursePrereqDataMap)
// 	if err != nil {
// 		t.Errorf("Error unmarshalling response body: %v", err)
// 	}
// 	coursePrereqDataMap2 := make(map[string]CoursePrerequisiteJson)
// 	err = json.Unmarshal(dataBytes, &coursePrereqDataMap2)
// 	if err != nil {
// 		t.Errorf("Error unmarshalling response body: %v", err)
// 	}

// 	for key := range coursePrereqDataMap {
// 		if _, ok := coursePrereqDataMap2[key]; !ok {
// 			t.Errorf("Error: key %v not found in coursePrereqDataMap2", key)
// 		}
// 	}

// 	// collect key and sort
// 	keys := make([]string, 0, len(coursePrereqDataMap2))
// 	for key := range coursePrereqDataMap2 {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)
// 	for _, key := range keys {
// 		fmt.Println(key, coursePrereqDataMap2[key].Prerequisites.TransformPrereq())
// 		// pq := coursePrereqDataMap2[key].Prerequisites.TransformPrereq()
// 		// fmt.Println(pq)
// 	}
// }
