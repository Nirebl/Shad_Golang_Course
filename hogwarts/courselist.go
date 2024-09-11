//go:build !solution

package hogwarts

import "sort"

func GetCourseList(prereqs map[string][]string) []string {
	visited := make(map[string]bool)
	onStack := make(map[string]bool)
	result := make([]string, 0)

	var hasCycle bool
	var dfs func(course string)
	dfs = func(course string) {
		if hasCycle {
			return
		}
		visited[course] = true
		onStack[course] = true
		for _, prereq := range prereqs[course] {
			if onStack[prereq] {
				hasCycle = true
				return
			}
			if !visited[prereq] {
				dfs(prereq)
			}
		}
		onStack[course] = false
		result = append(result, course)
	}

	courses := make([]string, 0, len(prereqs))
	for course := range prereqs {
		courses = append(courses, course)
	}
	sort.Strings(courses)

	for _, course := range courses {
		if !visited[course] {
			dfs(course)
		}
	}

	if hasCycle {
		panic("there is cycle")
	}
	return result
}
