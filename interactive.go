package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InteractiveModel struct {
	records    []DNSRecord
	filtered   []DNSRecord
	cf         *CloudflareManager
	cursor     int
	viewMode   string // "list", "add", "edit", "delete"
	page       int
	pageSize   int
	filter     string
	sortBy     string
	ascending  bool
}

func startInteractiveMode() {
	cf, err := NewCloudflareManager()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	records, err := cf.GetAllDNSRecords()
	if err != nil {
		fmt.Printf("获取DNS记录失败: %v\n", err)
		return
	}

	model := &InteractiveModel{
		records:    records,
		filtered:   records,
		cf:         cf,
		cursor:     0,
		viewMode:   "list",
		page:       0,
		pageSize:   20,
		filter:     "",
		sortBy:     "zone",
		ascending:  true,
	}

	// 启动交互式界面
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("启动交互式界面失败: %v\n", err)
		return
	}
}

func (m *InteractiveModel) Init() tea.Cmd {
	return nil
}

func (m *InteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.viewMode {
		case "list":
			return m.updateListMode(msg)
		case "add":
			return m.updateAddMode(msg)
		case "edit":
			return m.updateEditMode(msg)
		case "delete":
			return m.updateDeleteMode(msg)
		}
	}
	return m, nil
}

func (m *InteractiveModel) View() string {
	switch m.viewMode {
	case "list":
		return m.viewListMode()
	case "add":
		return m.viewAddMode()
	case "edit":
		return m.viewEditMode()
	case "delete":
		return m.viewDeleteMode()
	default:
		return m.viewListMode()
	}
}

func (m *InteractiveModel) updateListMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case "a":
		m.viewMode = "add"
	case "e":
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			m.viewMode = "edit"
		}
	case "d":
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			m.viewMode = "delete"
		}
	case "f":
		m.viewMode = "filter"
	case "s":
		m.ascending = !m.ascending
		m.applyFiltersAndSort()
	case "r":
		m.refreshRecords()
	}
	return m, nil
}

func (m *InteractiveModel) updateAddMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = "list"
	}
	return m, nil
}

func (m *InteractiveModel) updateEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = "list"
	}
	return m, nil
}

func (m *InteractiveModel) updateDeleteMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = "list"
	case "y":
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			record := m.filtered[m.cursor]
			err := m.cf.DeleteDNSRecord(record.ID, record.ZoneID)
			if err != nil {
				fmt.Printf("删除失败: %v\n", err)
			} else {
				fmt.Println("删除成功")
				m.refreshRecords()
			}
		}
		m.viewMode = "list"
	}
	return m, nil
}

func (m *InteractiveModel) viewListMode() string {
	var sb strings.Builder

	// 标题
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true).
		Render("🌐 Cloudflare DNS 管理器")
	sb.WriteString(title + "\n\n")

	// 状态栏
	status := fmt.Sprintf("总记录数: %d | 筛选后: %d | 排序: %s %s", 
		len(m.records), len(m.filtered), m.sortBy, 
		map[bool]string{true: "↑", false: "↓"}[m.ascending])
	sb.WriteString(lipgloss.NewStyle().Faint(true).Render(status) + "\n\n")

	// 操作提示
	help := "操作: ↑↓ 选择 | a 添加 | e 编辑 | d 删除 | f 筛选 | s 排序 | r 刷新 | q 退出"
	sb.WriteString(lipgloss.NewStyle().Faint(true).Render(help) + "\n\n")

	// 记录列表
	if len(m.filtered) == 0 {
		sb.WriteString("没有找到DNS记录\n")
		return sb.String()
	}

	// 表头
	header := fmt.Sprintf("%-20s %-15s %-8s %-40s %-6s %-8s",
		"域名", "记录名", "类型", "内容", "TTL", "代理")
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render(header) + "\n")
	sb.WriteString(strings.Repeat("-", 100) + "\n")

	// 分页显示
	start := m.page * m.pageSize
	end := start + m.pageSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		record := m.filtered[i]
		isSelected := i == m.cursor

		// 格式化内容
		content := record.Content
		if len(content) > 37 {
			content = content[:34] + "..."
		}

		proxiedStr := "否"
		if record.Proxied {
			proxiedStr = "是"
		}

		ttlStr := strconv.Itoa(record.TTL)
		if record.TTL == 1 {
			ttlStr = "自动"
		}

		line := fmt.Sprintf("%-20s %-15s %-8s %-40s %-6s %-8s",
			record.ZoneName, record.Name, record.Type, content, ttlStr, proxiedStr)

		if isSelected {
			sb.WriteString(lipgloss.NewStyle().
				Background(lipgloss.Color("#4A90E2")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Render("> " + line) + "\n")
		} else {
			sb.WriteString("  " + line + "\n")
		}
	}

	// 分页信息
	if len(m.filtered) > m.pageSize {
		totalPages := (len(m.filtered) + m.pageSize - 1) / m.pageSize
		pageInfo := fmt.Sprintf("\n第 %d/%d 页", m.page+1, totalPages)
		sb.WriteString(lipgloss.NewStyle().Faint(true).Render(pageInfo))
	}

	return sb.String()
}

func (m *InteractiveModel) viewAddMode() string {
	var sb strings.Builder
	sb.WriteString("添加DNS记录 (按ESC返回)\n\n")
	sb.WriteString("请使用命令行模式添加记录:\n")
	sb.WriteString("cf-dns-manager add [域名] [记录名] [类型] [内容] --ttl [TTL] --proxied\n\n")
	sb.WriteString("示例:\n")
	sb.WriteString("cf-dns-manager add example.com www A 192.168.1.1 --ttl 300 --proxied\n")
	return sb.String()
}

func (m *InteractiveModel) viewEditMode() string {
	var sb strings.Builder
	sb.WriteString("编辑DNS记录 (按ESC返回)\n\n")
	sb.WriteString("请使用命令行模式编辑记录:\n")
	sb.WriteString("cf-dns-manager update [记录ID] [域名] [记录名] [类型] [内容] --ttl [TTL] --proxied\n\n")
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		record := m.filtered[m.cursor]
		sb.WriteString(fmt.Sprintf("当前选中记录: %s\n", record.ID))
		sb.WriteString(fmt.Sprintf("域名: %s\n", record.ZoneName))
		sb.WriteString(fmt.Sprintf("记录名: %s\n", record.Name))
		sb.WriteString(fmt.Sprintf("类型: %s\n", record.Type))
		sb.WriteString(fmt.Sprintf("内容: %s\n", record.Content))
	}
	return sb.String()
}

func (m *InteractiveModel) viewDeleteMode() string {
	var sb strings.Builder
	sb.WriteString("删除DNS记录 (按ESC返回)\n\n")
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		record := m.filtered[m.cursor]
		sb.WriteString(fmt.Sprintf("确定要删除以下记录吗?\n\n"))
		sb.WriteString(fmt.Sprintf("记录ID: %s\n", record.ID))
		sb.WriteString(fmt.Sprintf("域名: %s\n", record.ZoneName))
		sb.WriteString(fmt.Sprintf("记录名: %s\n", record.Name))
		sb.WriteString(fmt.Sprintf("类型: %s\n", record.Type))
		sb.WriteString(fmt.Sprintf("内容: %s\n\n", record.Content))
		sb.WriteString("按 'y' 确认删除，按 'ESC' 取消")
	}
	return sb.String()
}

func (m *InteractiveModel) applyFiltersAndSort() {
	m.filtered = m.records

	// 应用筛选
	if m.filter != "" {
		filters := make(map[string]string)
		// 简单的关键词筛选
		for _, record := range m.records {
			if strings.Contains(strings.ToLower(record.Name), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(record.ZoneName), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(record.Content), strings.ToLower(m.filter)) {
				m.filtered = append(m.filtered, record)
			}
		}
	}

	// 应用排序
	SortRecords(m.filtered, m.sortBy, m.ascending)
}

func (m *InteractiveModel) refreshRecords() {
	records, err := m.cf.GetAllDNSRecords()
	if err != nil {
		return
	}
	m.records = records
	m.applyFiltersAndSort()
} 