package events

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var url_expression string = `https://(?:ptb\.|canary\.)?discord\.com/channels/(\d{17,19})/(\d{17,19})/(\d{17,19})`
var id_expression string = `(\d{17,19})/(\d{17,19})/(\d{17,19})`

func getIDs(url string) []string {
	url_regex, _ := regexp.Compile(url_expression)
	channels_regex, _ := regexp.Compile(id_expression)
	ids := strings.Split(channels_regex.FindAllString(url_regex.FindString(url), -1)[0], "/")
	return ids
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	url_regex, _ := regexp.Compile(url_expression)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if url_regex.MatchString(m.Content) {
		urls := url_regex.FindAllString(m.Content, -1);
		embeds := []*discordgo.MessageEmbed{}
		for _, url := range urls {
			ids := getIDs(url);
			message, err := s.State.Message(ids[1], ids[2])
			if err != nil {
				return
			}
			
			guild, err := s.State.Guild(message.GuildID)
			if err != nil {
				return
			}

			if len(message.Content) > 2048 {
				message.Content = message.Content[:2045] + "..."
			}
			timestamp, _ := discordgo.SnowflakeTimestamp(message.ID);

			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name: message.Author.Username, 
					IconURL: message.Author.AvatarURL("64"),
				},
				Description: message.Content,
				Timestamp: timestamp.Local().Format("2006-01-02 15:04:05"),
				Footer: &discordgo.MessageEmbedFooter{
					Text: guild.Name,
					IconURL: guild.IconURL("64"),
				},
			}

			if message.Attachments != nil {
				embed.Image = &discordgo.MessageEmbedImage{
					URL: message.Attachments[0].URL,
				}
			}

			embeds = append(embeds, embed)
		}
		s.ChannelMessageSendEmbedsReply(m.ChannelID, embeds, m.Reference());
	}
}

