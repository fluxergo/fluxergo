package fluxer

import (
	"fmt"
	"io"
	"slices"

	"github.com/disgoorg/snowflake/v2"
)

// NewMessageUpdate returns a new MessageUpdate with no fields set.
func NewMessageUpdate() MessageUpdate {
	return MessageUpdate{}
}

// MessageUpdate is used to edit a Message.
type MessageUpdate struct {
	Content         *string             `json:"content,omitempty"`
	Embeds          *[]Embed            `json:"embeds,omitempty"`
	Attachments     *[]AttachmentUpdate `json:"attachments,omitempty"`
	Files           []*File             `json:"-"`
	AllowedMentions *AllowedMentions    `json:"allowed_mentions,omitempty"`
	// Flags are the MessageFlags of the message.
	// Be careful not to override the current flags when editing messages from other users - this will result in a permission error.
	// Use MessageFlags.Add for flags like fluxer.MessageFlagSuppressEmbeds.
	Flags *MessageFlags `json:"flags,omitempty"`
}

func (MessageUpdate) interactionCallbackData() {}

// ToBody returns the MessageUpdate ready for body.
func (m MessageUpdate) ToBody() (any, error) {
	if len(m.Files) > 0 {
		for _, attachmentCreate := range parseAttachments(m.Files) {
			if m.Attachments == nil {
				m.Attachments = new([]AttachmentUpdate)
			}
			*m.Attachments = append(*m.Attachments, attachmentCreate)
		}
		return PayloadWithFiles(m, m.Files...)
	}
	return m, nil
}

// WithContent returns a new MessageUpdate with the provided content.
func (m MessageUpdate) WithContent(content string) MessageUpdate {
	m.Content = &content
	return m
}

// WithContentf returns a new MessageUpdate with the formatted content.
func (m MessageUpdate) WithContentf(content string, a ...any) MessageUpdate {
	return m.WithContent(fmt.Sprintf(content, a...))
}

// ClearContent returns a new MessageUpdate with no content.
func (m MessageUpdate) ClearContent() MessageUpdate {
	return m.WithContent("")
}

// WithEmbeds returns a new MessageUpdate with the provided Embed(s).
func (m MessageUpdate) WithEmbeds(embeds ...Embed) MessageUpdate {
	m.Embeds = &embeds
	return m
}

// WithEmbed returns a new MessageUpdate with the provided Embed at the index.
func (m MessageUpdate) WithEmbed(i int, embed Embed) MessageUpdate {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	if len(*m.Embeds) > i {
		newEmbeds := slices.Insert(*m.Embeds, i, embed)
		m.Embeds = &newEmbeds
	}
	return m
}

// AddEmbeds returns a new MessageUpdate with the provided embeds added.
func (m MessageUpdate) AddEmbeds(embeds ...Embed) MessageUpdate {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	newEmbeds := append(slices.Clone(*m.Embeds), embeds...)
	m.Embeds = &newEmbeds
	return m
}

// ClearEmbeds returns a new MessageUpdate with no embeds.
func (m MessageUpdate) ClearEmbeds() MessageUpdate {
	m.Embeds = &[]Embed{}
	return m
}

// RemoveEmbed returns a new MessageUpdate with the embed at the index removed.
func (m MessageUpdate) RemoveEmbed(i int) MessageUpdate {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	if len(*m.Embeds) > i {
		newEmbeds := slices.Delete(slices.Clone(*m.Embeds), i, i+1)
		m.Embeds = &newEmbeds
	}
	return m
}

// WithFiles returns a new MessageUpdate with the provided File(s).
func (m MessageUpdate) WithFiles(files ...*File) MessageUpdate {
	m.Files = files
	return m
}

// UpdateFile returns a new MessageUpdate with the File at the index.
func (m MessageUpdate) UpdateFile(i int, file *File) MessageUpdate {
	if len(m.Files) > i {
		m.Files = slices.Clone(m.Files)
		m.Files[i] = file
	}
	return m
}

// AddFiles returns a new MessageUpdate with the File(s) added.
func (m MessageUpdate) AddFiles(files ...*File) MessageUpdate {
	m.Files = append(m.Files, files...)
	return m
}

// AddFile returns a new MessageUpdate with a File added.
func (m MessageUpdate) AddFile(name string, description string, reader io.Reader, flags ...FileFlags) MessageUpdate {
	m.Files = append(m.Files, NewFile(name, description, reader, flags...))
	return m
}

// ClearFiles returns a new MessageUpdate with no File(s).
func (m MessageUpdate) ClearFiles() MessageUpdate {
	m.Files = []*File{}
	return m
}

// RemoveFile returns a new MessageUpdate with the File at the index removed.
func (m MessageUpdate) RemoveFile(i int) MessageUpdate {
	if len(m.Files) > i {
		m.Files = slices.Delete(slices.Clone(m.Files), i, i+1)
	}
	return m
}

// RetainAttachments returns a new MessageUpdate that retains the provided Attachment(s).
func (m MessageUpdate) RetainAttachments(attachments ...Attachment) MessageUpdate {
	if m.Attachments == nil {
		m.Attachments = &[]AttachmentUpdate{}
	}
	newAttachments := slices.Clone(*m.Attachments)
	for _, attachment := range attachments {
		newAttachments = append(newAttachments, AttachmentKeep{
			ID: attachment.ID,
		})
	}
	m.Attachments = &newAttachments
	return m
}

// RetainAttachmentsByID returns a new MessageUpdate that retains the Attachment(s) with the provided IDs.
func (m MessageUpdate) RetainAttachmentsByID(attachmentIDs ...snowflake.ID) MessageUpdate {
	if m.Attachments == nil {
		m.Attachments = &[]AttachmentUpdate{}
	}
	newAttachments := slices.Clone(*m.Attachments)
	for _, attachmentID := range attachmentIDs {
		newAttachments = append(newAttachments, AttachmentKeep{
			ID: attachmentID,
		})
	}
	m.Attachments = &newAttachments
	return m
}

// WithAllowedMentions returns a new MessageUpdate with the provided AllowedMentions.
func (m MessageUpdate) WithAllowedMentions(allowedMentions *AllowedMentions) MessageUpdate {
	m.AllowedMentions = allowedMentions
	return m
}

// ClearAllowedMentions returns a new MessageUpdate with no AllowedMentions.
func (m MessageUpdate) ClearAllowedMentions() MessageUpdate {
	return m.WithAllowedMentions(nil)
}

// WithFlags returns a new MessageUpdate with the provided MessageFlags.
// Be careful not to override the current flags when editing messages from other users - this will result in a permission error.
// Use WithSuppressEmbeds or AddFlags for flags like MessageFlagSuppressEmbeds.
func (m MessageUpdate) WithFlags(flags MessageFlags) MessageUpdate {
	m.Flags = &flags
	return m
}

// AddFlags returns a new MessageUpdate with the provided MessageFlags added.
func (m MessageUpdate) AddFlags(flags ...MessageFlags) MessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	newFlags := (*m.Flags).Add(flags...)
	m.Flags = &newFlags
	return m
}

// RemoveFlags returns a new MessageUpdate with the provided MessageFlags removed.
func (m MessageUpdate) RemoveFlags(flags ...MessageFlags) MessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	newFlags := (*m.Flags).Remove(flags...)
	m.Flags = &newFlags
	return m
}

// ClearFlags returns a new MessageUpdate with no MessageFlags.
func (m MessageUpdate) ClearFlags() MessageUpdate {
	return m.WithFlags(MessageFlagsNone)
}

// WithSuppressEmbeds returns a new MessageUpdate with MessageFlagSuppressEmbeds added/removed.
func (m MessageUpdate) WithSuppressEmbeds(suppressEmbeds bool) MessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	flags := *m.Flags
	if suppressEmbeds {
		flags = m.Flags.Add(MessageFlagSuppressEmbeds)
	} else {
		flags = m.Flags.Remove(MessageFlagSuppressEmbeds)
	}
	m.Flags = &flags
	return m
}

// WithIsComponentsV2 returns a new MessageUpdate with MessageFlagIsComponentsV2 added/removed.
// Once a message with the flag has been sent, it cannot be removed by editing the message.
func (m MessageUpdate) WithIsComponentsV2(isComponentV2 bool) MessageUpdate {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	flags := *m.Flags
	if isComponentV2 {
		flags = m.Flags.Add(MessageFlagIsComponentsV2)
	} else {
		flags = m.Flags.Remove(MessageFlagIsComponentsV2)
	}
	m.Flags = &flags
	return m
}
