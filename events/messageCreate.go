package events

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"log/slog"
)

var discord_url_expression string = `https://(?:ptb\.|canary\.)?discord\.com/channels/(\d{17,19})/(\d{17,19})/(\d{17,19})`
var url_expression string = `https?://[\w/:%#\$&\?\(\)~\.=\+\-]+`
var id_expression = `(\d{17,19})/(\d{17,19})/(\d{17,19})`

var discord_url_regex = regexp.MustCompile(discord_url_expression)
var url_regex = regexp.MustCompile(url_expression)
var ids_regex = regexp.MustCompile(id_expression)

type messageInfo struct {
	guild string;
	channel string;
	message string;
}

type discordInfo struct {
	guild *discordgo.Guild;
	channel *discordgo.Channel;
	message *discordgo.Message;
}

func getIDs(url string) ( messageInfo ) {
	ids := strings.Split(ids_regex.FindAllString(url_regex.FindString(url), -1)[0], "/")
	return messageInfo{guild: ids[0], channel: ids[1], message: ids[2]}
}

func getDiscordInformationFromURL(s *discordgo.Session, url string) ( discordInfo, error ) {
	msg_info := getIDs(url);
	guild, err := s.State.Guild(msg_info.guild)
	if err != nil {
		slog.Error("An error occured when retrieving guild", err.Error())
		return discordInfo{}, err
	}
	channel, err := s.Channel(msg_info.channel)
	if err != nil {
		slog.Error("An error occured when retrieving channel", err.Error())
		return discordInfo{}, err
	}
	message, err := s.ChannelMessage(msg_info.channel, msg_info.message)
	if err != nil {
		slog.Error("An error occured when retrieving message", err.Error())
		return discordInfo{}, err
	}

	return discordInfo{guild: guild, channel: channel, message: message}, nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		slog.Info("Ignoring bot message")
	}

	if discord_url_regex.MatchString(m.Content) {
		urls := discord_url_regex.FindAllString(m.Content, -1);
		embeds := []*discordgo.MessageEmbed{}
		for _, url := range urls {
		
			discord_info, err := getDiscordInformationFromURL(s, url)
			if err != nil {
				return
			}

			message := discord_info.message
			guild := discord_info.guild

			if url_regex.MatchString(message.Content) {
				for _, url := range url_regex.FindAllString(message.Content, -1) {
					if discord_url_regex.MatchString(url) {
						discord_info, err := getDiscordInformationFromURL(s, url)
						if err != nil {
							return
						}
						message.Content = strings.ReplaceAll(message.Content, url, "[" + discord_info.message.Content[8:28] + "](" + url + ")")

					} else {
						message.Content = strings.ReplaceAll(message.Content, url, "[" + url[8:28] + "](" + url + ")")
					}
				}
			}

			if len(message.Content) > 2048 {
				message.Content = message.Content[:2045] + "..."
				slog.Info("The message content is too long, it has been truncated")
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

			if len(message.Attachments) != 0 {
				embed.Image = &discordgo.MessageEmbedImage{
					URL: message.Attachments[0].URL,
				}
			}

			embeds = append(embeds, embed)
		}
		_,e := s.ChannelMessageSendEmbedsReply(m.ChannelID, embeds, m.Reference());
		if e != nil {
			println(e.Error())
		}
	}
}

