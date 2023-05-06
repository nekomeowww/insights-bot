package tgbot

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf16"

	"github.com/Junzki/link-preview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"

	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/nekomeowww/insights-bot/pkg/utils"
)

var (
	escapeForMarkdownV2MarkdownLinkRegexp = regexp.MustCompile(`(\[[^\][]*]\(http[^()]*\))|[_*[\]()~>#+=|{}.!-]`)
)

// EscapeForMarkdownV2
//
// 1. 任何字符码表在 1 到 126 之间的字符都可以加前缀 '\' 字符来转义，在这种情况下，它被视为一个普通字符，而不是标记的一部分。这意味着 '\' 字符通常必须加前缀 '\' 字符来转义。
// 2. 在 pre 和 code 的实体中，所有 '`' 和 '\' 字符都必须加前缀 '\' 字符转义。
// 3. 在所有其他地方，字符 '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!' 必须加前缀字符 '\' 转义。
//
// https://core.telegram.org/bots/api#formatting-options
func EscapeStringForMarkdownV2(src string) string {
	var result string

	escapingIndexes := make([][]int, 0)

	// 查询需要转义的字符
	for _, match := range escapeForMarkdownV2MarkdownLinkRegexp.FindAllStringSubmatchIndex(src, -1) {
		if match[2] != -1 && match[3] != -1 { // 匹配到了 markdown 链接
			continue // 不需要转义
		}

		escapingIndexes = append(escapingIndexes, match) // 需要转义
	}

	// 对需要转义的字符进行转义
	var lastMatchedIndex int

	for i, match := range escapingIndexes {
		if i == 0 {
			result += src[lastMatchedIndex:match[0]]
		} else {
			result += src[escapingIndexes[i-1][1]:match[0]]
		}

		result += `\` + src[match[0]:match[1]]
		lastMatchedIndex = match[1]
	}
	if lastMatchedIndex < len(src) {
		result += src[lastMatchedIndex:]
	}

	return result
}

func NewCallbackQueryData(component string, route string, queries url.Values) string {
	return fmt.Sprintf("cbq://%s/%s?%s", component, route, queries.Encode())
}

func FullNameFromFirstAndLastName(firstName, lastName string) string {
	if lastName == "" {
		return firstName
	}
	if firstName == "" {
		return lastName
	}
	if utils.ContainsCJKChar(firstName) && !utils.ContainsCJKChar(lastName) {
		return firstName + " " + lastName
	}
	if !utils.ContainsCJKChar(firstName) && utils.ContainsCJKChar(lastName) {
		return lastName + " " + firstName
	}
	if utils.ContainsCJKChar(firstName) && utils.ContainsCJKChar(lastName) {
		return lastName + " " + firstName
	}

	return firstName + " " + lastName
}

func ExtractTextFromMessage(message *tgbotapi.Message) string {
	text := lo.Ternary(message.Caption != "", message.Caption, message.Text)
	type MarkdownLink struct {
		Markdown []uint16
		Start    int
		End      int
	}

	text_utf16 := utf16.Encode([]rune(text))
	links := lop.Map(message.Entities, func(entity tgbotapi.MessageEntity, i int) MarkdownLink {
		start_i := entity.Offset
		end_i := start_i + entity.Length
		var title string
		var href string
		if entity.Type == "url" {
			href = string(utf16.Decode(text_utf16[start_i:end_i]))
			result, err := LinkPreview.PreviewLink(href, nil)
			if err != nil {
				return MarkdownLink{[]uint16{}, -1, -1}
			}
			title = result.Title
		} else if entity.Type == "text_link" {
			title = string(utf16.Decode(text_utf16[start_i:end_i]))
			href = entity.URL
		} else {
			return MarkdownLink{[]uint16{}, -1, -1}
		}
		unescaped, err := url.QueryUnescape(href)
		if err == nil {
			href = strings.ReplaceAll(unescaped, " ", "+")
		}
		md := "[" + title + "](" + href + ")"
		md_utf16 := utf16.Encode([]rune(md))
		return MarkdownLink{md_utf16, start_i, end_i}
	})

	for i := len(links) - 1; i >= 0; i-- {
		if links[i].Start == -1 {
			continue
		}
		text_utf16 = append(text_utf16[:links[i].Start], append(links[i].Markdown, text_utf16[links[i].End:]...)...)
	}

	return string(utf16.Decode(text_utf16))
}

// EscapeHTMLSymbols
//
//	< with &lt;
//	> with &gt;
//	& with &amp;
func EscapeHTMLSymbols(str string) string {
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	str = strings.ReplaceAll(str, "&", "&amp;")

	return str
}

var (
	matchMdTitles = regexp.MustCompile(`(?m)^(#){1,6} (.)*(\n)?`)
)

func ReplaceMarkdownTitlesToTelegramBoldElement(text string) (string, error) {
	return matchMdTitles.ReplaceAllStringFunc(text, func(s string) string {
		// remove hashtag
		for strings.HasPrefix(s, "#") {
			s = strings.TrimPrefix(s, "#")
		}
		// remove space
		s = strings.TrimPrefix(s, " ")

		sRunes := []rune(s)
		ret := "<b>" + string(sRunes[:len(sRunes)-1])

		// if the line ends with a newline, add a newline to the end of the bold element
		if strings.HasSuffix(s, "\n") {
			return ret + "</b>\n"
		}

		// otherwise, just return the bold element
		return ret + string(sRunes[len(sRunes)-1]) + "</b>"
	}), nil
}

func MapChatTypeToChineseText(chatType telegram.ChatType) string {
	switch chatType {
	case telegram.ChatTypePrivate:
		return "私聊"
	case telegram.ChatTypeGroup:
		return "群组"
	case telegram.ChatTypeSuperGroup:
		return "超级群组"
	case telegram.ChatTypeChannel:
		return "频道"
	default:
		return "未知"
	}
}

func MapMemberStatusToChineseText(memberStatus telegram.MemberStatus) string {
	switch memberStatus {
	case telegram.MemberStatusCreator:
		return "创建者"
	case telegram.MemberStatusAdministrator:
		return "管理员"
	case telegram.MemberStatusMember:
		return "成员"
	case telegram.MemberStatusRestricted:
		return "受限成员"
	case telegram.MemberStatusLeft:
		return "已离开"
	case telegram.MemberStatusKicked:
		return "已被踢出"
	default:
		return "未知"
	}
}
