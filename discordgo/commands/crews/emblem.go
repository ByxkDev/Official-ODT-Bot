package crews

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"

	_ "image/jpeg"
	_ "golang.org/x/image/webp"
)

var (
	MaxFileSize          int64 = 64 * 1024
	AllowedTypes               = []string{".png", ".jpeg", ".webp", ".jpg"}
	EmblemDir                  = "C:/Users/Administrator/Desktop/ODT-Server-main/Storage/Crews"
	CompressonatorBinary       = "C:/Compressonator/bin/CLI/compressonatorcli.exe"
)

func isAllowed(ext string) bool {
	ext = strings.ToLower(ext)

	for _, t := range AllowedTypes {
		if ext == t {
			return true
		}
	}

	return false
}

func Emblem(s *discordgo.Session, m *discordgo.MessageCreate) {

	if len(m.Attachments) == 0 {
		utils.SendError(s, m, "Missing File", "Usage: ?crewemblem <attach image>")
		return
	}

	attachment := m.Attachments[0]
	ext := strings.ToLower(filepath.Ext(attachment.Filename))

	if !isAllowed(ext) {
		utils.SendError(s, m, "Invalid File", "Allowed: PNG, JPG, WEBP under 64KB")
		return
	}

	if int64(attachment.Size) > MaxFileSize {
		utils.SendError(s, m, "File Too Large", "Max size is 64KB.")
		return
	}

	var crewID string

	err := utils.DB().QueryRow("SELECT crew_id FROM crews WHERE crew_owner = ?", m.Author.ID).Scan(&crewID)
	if err != nil || crewID == "" {
		utils.SendError(s, m, "No Crew Found", "You must be a crew owner.")
		return
	}

	emblemDir := filepath.Join(EmblemDir, crewID)
	tempPNG := filepath.Join(emblemDir, "emblem_128.png")
	finalDDS := filepath.Join(emblemDir, "emblem_128.dds")

	err = os.MkdirAll(emblemDir, 0755)
	if err != nil {
		utils.SendError(s, m, "Error", "Failed creating emblem folder.")
		return
	}

	resp, err := http.Get(attachment.URL)
	if err != nil {
		utils.SendError(s, m, "Download Failed", "Could not download attachment.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.SendError(s, m, "Download Failed", "Invalid Discord response.")
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SendError(s, m, "Error", "Failed reading image.")
		return
	}

	if int64(len(data)) > MaxFileSize {
		utils.SendError(s, m, "File Too Large", "Max size is 64KB.")
		return
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		utils.SendError(s, m, "Invalid Image", "Could not decode image.")
		return
	}

	fmt.Println("Uploaded image format:", format)

	file, err := os.Create(tempPNG)
	if err != nil {
		utils.SendError(s, m, "Error", "Could not create PNG file.")
		return
	}

	err = png.Encode(file, img)
	file.Close()

	if err != nil {
		utils.SendError(s, m, "Error", "Failed converting image.")
		return
	}

	cmd := exec.Command(CompressonatorBinary, "-fd", "BC3", tempPNG, finalDDS)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err = cmd.Run()
	if err != nil {
		fmt.Println("[COMPRESSION ERROR]", output.String())
		utils.SendError(s, m, "Compression Failed", output.String())
		return
	}

	_ = os.Remove(tempPNG)
	files, _ := os.ReadDir(emblemDir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "emblem_128") && f.Name() != "emblem_128.dds" {
			_ = os.Remove(filepath.Join(emblemDir, f.Name()))
		}
	}

	utils.SendWebhook(s, fmt.Sprintf("New Crew Emblem Uploaded\nUpdated by **%s**", m.Author.Username))
	utils.SendEmbed(s, m, "Crew Emblem Updated", "Your emblem has been successfully updated.")
}