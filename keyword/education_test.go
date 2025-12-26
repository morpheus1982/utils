package keyword

import (
	"reflect"
	"testing"
)

func TestExtractEducationKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *EducationKeywords
	}{
		{
			name:  "学校类关键字提取",
			input: "我的孩子在成都七中上学，同时也考虑石室中学",
			expected: &EducationKeywords{
				Schools: []string{"七中", "石室中学", "石室", "成都七中"},
				Grades:  []string{},
				Exams:   []string{},
				Careers: []string{},
			},
		},
		{
			name:  "年级类关键字提取",
			input: "孩子正在上小学三年级，准备升入初中一年级",
			expected: &EducationKeywords{
				Schools: []string{},
				Grades:  []string{"小学", "一年级", "三年级", "初中"},
				Exams:   []string{},
				Careers: []string{},
			},
		},
		{
			name:  "考试类关键字提取",
			input: "今年要参加中考和高考，还要准备司法考试和一建考试",
			expected: &EducationKeywords{
				Schools: []string{},
				Grades:  []string{},
				Exams:   []string{"中考", "高考", "司法考试", "法考", "司法", "一建"},
				Careers: []string{},
			},
		},
		{
			name:  "职业类关键字提取",
			input: "我想成为一名教师或者律师，也可以考虑医生这个职业",
			expected: &EducationKeywords{
				Schools: []string{},
				Grades:  []string{},
				Exams:   []string{},
				Careers: []string{"教师", "医生", "律师"},
			},
		},
		{
			name:  "混合关键字提取",
			input: "成都七中的高三学生正在准备高考，希望将来成为一名教师",
			expected: &EducationKeywords{
				Schools: []string{"成都七中", "七中"},
				Grades:  []string{"高三"},
				Exams:   []string{"高考"},
				Careers: []string{"教师"},
			},
		},
		{
			name:     "无关键字文本",
			input:    "今天天气很好，我们去公园散步",
			expected: nil,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractEducationKeywords(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractEducationKeywords(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractSchoolKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "提取学校全称和简称",
			input:    "我的孩子在成都七中上学",
			expected: []string{"成都七中", "七中"},
		},
		{
			name:     "提取多个学校",
			input:    "考虑成都七中和石室中学",
			expected: []string{"成都七中", "七中", "石室中学", "石室"},
		},
		{
			name:     "无学校关键字",
			input:    "今天天气很好",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSchoolKeywords(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractSchoolKeywords(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractGradeKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "提取小学年级",
			input:    "孩子正在上小学三年级",
			expected: []string{"小学", "三年级"},
		},
		{
			name:     "提取多个年级",
			input:    "从小学一年级到初中三年级",
			expected: []string{"小学", "一年级", "初中", "三年级"},
		},
		{
			name:     "无年级关键字",
			input:    "今天天气很好",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractGradeKeywords(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractGradeKeywords(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractExamKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "提取升学考试",
			input:    "今年要参加中考和高考",
			expected: []string{"中考", "高考"},
		},
		{
			name:     "提取资格考试",
			input:    "准备司法考试和一建考试",
			expected: []string{"司法考试", "法考", "司法", "一建"},
		},
		{
			name:     "无考试关键字",
			input:    "今天天气很好",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractExamKeywords(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractExamKeywords(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractCareerKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "提取教育职业",
			input:    "我想成为一名教师",
			expected: []string{"教师"},
		},
		{
			name:     "提取多个职业",
			input:    "考虑教师或者律师职业",
			expected: []string{"教师", "律师"},
		},
		{
			name:     "无职业关键字",
			input:    "今天天气很好",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCareerKeywords(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractCareerKeywords(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
