package keyword

import (
	"strings"
)

// EducationKeywords 教育相关关键字提取结果
type EducationKeywords struct {
	Schools []string // 学校类关键字
	Grades  []string // 年级类关键字
	Exams   []string // 考试类关键字
	Careers []string // 职业类关键字
}

// 学校类关键字映射表 - 全称到简称的映射
var schoolKeywordMap = map[string][]string{
	// 成都地区知名学校
	"成都七中":      {"七中", "成都七中", "成都七中育才", "七中育才"},
	"成都七中育才学校":  {"七中育才", "七中", "成都七中育才"},
	"石室中学":      {"石室", "石室中学", "成都四中", "四中"},
	"成都石室中学":    {"石室", "石室中学", "成都四中", "四中"},
	"树德中学":      {"树德", "树德中学", "成都九中", "九中"},
	"成都树德中学":    {"树德", "树德中学", "成都九中", "九中"},
	"成都外国语学校":   {"成外", "成都外国语", "成外"},
	"成都实验外国语学校": {"实外", "实验外国语", "实外"},
	"成都嘉祥外国语学校": {"嘉祥", "嘉祥外国语", "嘉祥"},
	"西川中学":      {"西川", "西川中学"},
	"师大一中":      {"师大一中", "师大一中"},

	// 北京地区知名学校
	"人大附中":    {"人大附", "人大附中", "RDFZ"},
	"北京四中":    {"四中", "北京四中"},
	"北京八中":    {"八中", "北京八中"},
	"北师大附中":   {"师大附中", "北师大附中"},
	"清华附中":    {"清华附", "清华附中"},
	"北大附中":    {"北大附", "北大附中"},
	"北京101中学": {"101中学", "101中"},

	// 上海地区知名学校
	"上海中学":    {"上海中学", "上中"},
	"华东师大二附中": {"华师大二附中", "华二"},
	"复旦附中":    {"复旦附", "复旦附中"},
	"交大附中":    {"交大附", "交大附中"},

	// 其他地区知名学校
	"黄冈中学":    {"黄冈中学", "黄冈"},
	"衡水中学":    {"衡水中学", "衡水"},
	"南京外国语学校": {"南外", "南京外国语"},
	"杭州第二中学":  {"杭二中", "杭州二中"},
	"深圳中学":    {"深圳中学", "深中"},
	"青桐鸣":     {"青桐鸣", "青桐鸣教育"},
	"新东方":     {"新东方", "新东方教育"},
	"学而思":     {"学而思", "学而思教育"},
	"好未来":     {"好未来", "好未来教育"},
}

// 年级类关键字
var gradeKeywords = []string{
	// 学前阶段
	"幼儿园", "学前班", "幼小衔接", "幼小",

	// 小学阶段
	"小学", "一年级", "二年级", "三年级", "四年级", "五年级", "六年级",
	"1年级", "2年级", "3年级", "4年级", "5年级", "6年级",
	"小一", "小二", "小三", "小四", "小五", "小六",

	// 初中阶段
	"初中", "初一", "初二", "初三",
	"七年级", "八年级", "九年级",
	"7年级", "8年级", "9年级",
	"初一期末", "初二期末", "初三期末",

	// 高中阶段
	"高中", "高一", "高二", "高三",
	"十年级", "十一年级", "十二年级",
	"10年级", "11年级", "12年级",
	"高一期末", "高二期末", "高三期末",

	// 其他阶段
	"小升初", "初升高",
}

// 考试类关键字
var examKeywords = []string{
	// 国内升学考试
	"中考", "高考", "小升初", "初升高",
	"联考", "统考", "模拟考试", "月考", "期中考试", "期末考试",
	"一模", "二模", "三模",

	// 专业资格考试
	"司法考试", "法考", "司法",
	"一级建造师", "一建", "二级建造师", "二建",
	"注册会计师", "注会", "CPA", "ACCA",
	"教师资格证", "教资", "教师证",
	"医师资格考试", "医师资格证", "医师",
	"护士执业资格证", "护士证", "护士",
	"执业药师", "药师证", "药师",
	"注册税务师", "税务师", "CTA",
	"注册造价师", "造价师", "造价",
	"注册安全工程师", "安全工程师", "安全工程师",

	// 外语考试
	"英语四级", "四级", "CET4", "英语六级", "六级", "CET6",
	"雅思", "IELTS", "托福", "TOEFL", "GRE", "GMAT",
	"日语等级考试", "JLPT", "N1", "N2", "N3", "N4", "N5",
	"韩语等级考试", "TOPIK", "德福", "TestDaF",
	"法语水平考试", "TEF", "TCF",

	// 公务员考试
	"公务员考试", "国考", "省考", "公务员",
	"事业单位考试", "事业单位", "事业编",
	"选调生考试", "选调生",

	// 企业入职考试
	"银行入职考试", "银行考试", "银行招聘",
	"国家电网考试", "电网考试", "电网招聘",
	"中石油考试", "中石油招聘", "中石化考试", "中石化招聘",
	"中国移动考试", "移动招聘", "中国联通考试", "联通招聘",
	"中国电信考试", "电信招聘",
	"腾讯招聘", "阿里招聘", "百度招聘", "字节跳动招聘",
	"华为招聘", "小米招聘", "京东招聘", "美团招聘",

	// 其他考试
	"考研", "研究生入学考试", "保研",
	"专升本", "专插本", "专接本",
	"成人高考", "自考", "网络教育",
}

// 职业类关键字
var careerKeywords = []string{
	// 教育类职业
	"教师", "老师", "教授", "讲师", "助教", "班主任",
	"校长", "院长", "系主任", "辅导员",

	// 医疗类职业
	"医生", "医师", "护士", "护师", "药师", "药剂师",
	"牙医", "兽医", "心理咨询师", "心理医生",

	// 法律类职业
	"律师", "法官", "检察官", "法警", "法律顾问",

	// 金融类职业
	"会计师", "注册会计师", "审计师", "精算师",
	"银行家", "投资顾问", "理财师", "金融分析师",

	// 技术类职业
	"工程师", "软件工程师", "硬件工程师", "网络工程师",
	"建筑师", "设计师", "产品经理", "项目经理",

	// 管理类职业
	"经理", "总监", "CEO", "CTO", "CFO", "COO",
	"总经理", "副总经理", "部门经理", "主管",

	// 其他职业
	"记者", "编辑", "作家", "翻译", "导游", "厨师",
	"警察", "消防员", "军人", "士兵", "军官",
}

// ExtractEducationKeywords 从文本中提取教育相关关键字
func ExtractEducationKeywords(text string) *EducationKeywords {
	// 清理文本，去除多余空格和特殊字符
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	// 初始化结果
	result := &EducationKeywords{}

	// 提取学校类关键字
	result.Schools = extractSchoolKeywords(text)

	// 提取年级类关键字
	result.Grades = extractGradeKeywords(text)

	// 提取考试类关键字
	result.Exams = extractExamKeywords(text)

	// 提取职业类关键字
	result.Careers = extractCareerKeywords(text)

	// 如果所有字段都为空，返回nil
	if len(result.Schools) == 0 && len(result.Grades) == 0 && len(result.Exams) == 0 && len(result.Careers) == 0 {
		return nil
	}

	return result
}

// extractSchoolKeywords 从文本中提取学校类关键字
func extractSchoolKeywords(text string) []string {
	var results []string
	found := make(map[string]bool)

	// 首先检查全称
	for fullName, abbreviations := range schoolKeywordMap {
		if strings.Contains(text, fullName) {
			if !found[fullName] {
				results = append(results, fullName)
				found[fullName] = true
			}
		}

		// 检查简称
		for _, abbr := range abbreviations {
			if strings.Contains(text, abbr) && !found[abbr] {
				results = append(results, abbr)
				found[abbr] = true
			}
		}
	}

	return results
}

// extractGradeKeywords 从文本中提取年级类关键字
func extractGradeKeywords(text string) []string {
	var results []string
	found := make(map[string]bool)

	for _, keyword := range gradeKeywords {
		if strings.Contains(text, keyword) && !found[keyword] {
			results = append(results, keyword)
			found[keyword] = true
		}
	}

	return results
}

// extractExamKeywords 从文本中提取考试类关键字
func extractExamKeywords(text string) []string {
	var results []string
	found := make(map[string]bool)

	for _, keyword := range examKeywords {
		if strings.Contains(text, keyword) && !found[keyword] {
			results = append(results, keyword)
			found[keyword] = true
		}
	}

	return results
}

// extractCareerKeywords 从文本中提取职业类关键字
func extractCareerKeywords(text string) []string {
	var results []string
	found := make(map[string]bool)

	for _, keyword := range careerKeywords {
		if strings.Contains(text, keyword) && !found[keyword] {
			results = append(results, keyword)
			found[keyword] = true
		}
	}

	return results
}

// ContainsEducationKeywords 检查文本是否包含教育相关关键字
func ContainsEducationKeywords(text string) bool {
	keywords := ExtractEducationKeywords(text)
	return keywords != nil
}

// CountEducationKeywords 统计文本中教育相关关键字的种类数量
func CountEducationKeywords(text string) int {
	keywords := ExtractEducationKeywords(text)
	if keywords == nil {
		return 0
	}

	count := 0
	if len(keywords.Schools) > 0 {
		count++
	}
	if len(keywords.Grades) > 0 {
		count++
	}
	if len(keywords.Exams) > 0 {
		count++
	}
	if len(keywords.Careers) > 0 {
		count++
	}

	return count
}
