package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "toolbox",
		Short: "Toolbox is a CLI tool with a GUI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var guiCmd = &cobra.Command{
		Use:   "gui",
		Short: "Launch the GUI",
		Run: func(cmd *cobra.Command, args []string) {
			launchGUI()
		},
	}

	var mediaRestartCmd = &cobra.Command{
		Use:   "mediaRestart",
		Short: "Restart media services",
		Run: func(cmd *cobra.Command, args []string) {
			mediaRestart()
		},
	}

	var musicStartCmd = &cobra.Command{
		Use:   "musicStart",
		Short: "Start music playback",
		Run: func(cmd *cobra.Command, args []string) {
			musicStart()
		},
	}

	var fileMvUpCmd = &cobra.Command{
		Use:   "fileMvUp",
		Short: "Move files up from subdirectories to the specified directory",
		Run: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			fileMvUp(path)
		},
	}
	fileMvUpCmd.Flags().StringP("path", "p", "", "Path to the directory")
	fileMvUpCmd.MarkFlagRequired("path")

	var ipListCmd = &cobra.Command{
		Use:   "ipList",
		Short: "List IP addresses",
		Run: func(cmd *cobra.Command, args []string) {
			ipList()
		},
	}

	var fishEditCmd = &cobra.Command{
		Use:   "fishEdit",
		Short: "Edit the fish shell configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fishEdit()
		},
	}

	rootCmd.AddCommand(fishEditCmd)

	rootCmd.AddCommand(ipListCmd)
	rootCmd.AddCommand(guiCmd)
	rootCmd.AddCommand(mediaRestartCmd)
	rootCmd.AddCommand(musicStartCmd)
	rootCmd.AddCommand(fileMvUpCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func launchGUI() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Toolbox GUI")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	vbox.SetMarginStart(20)
	vbox.SetMarginEnd(20)
	vbox.SetMarginTop(20)
	vbox.SetMarginBottom(20)

	mediaRestartButton, err := gtk.ButtonNewWithLabel("Restart Media")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	mediaRestartButton.SetMarginTop(5)
	mediaRestartButton.SetMarginBottom(5)
	mediaRestartButton.SetSizeRequest(150, -1)
	mediaRestartButton.Connect("clicked", func() {
		mediaRestart()
	})
	vbox.PackStart(mediaRestartButton, false, false, 0)

	musicStartButton, err := gtk.ButtonNewWithLabel("Start Music")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	musicStartButton.SetMarginTop(5)
	musicStartButton.SetMarginBottom(5)
	musicStartButton.SetSizeRequest(150, -1)
	musicStartButton.Connect("clicked", func() {
		musicStart()
	})
	vbox.PackStart(musicStartButton, false, false, 0)

	musicEndButton, err := gtk.ButtonNewWithLabel("Stop Music")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	musicEndButton.SetMarginTop(5)
	musicEndButton.SetMarginBottom(5)
	musicEndButton.SetSizeRequest(150, -1)
	musicEndButton.Connect("clicked", func() {
		musicEnd()
	})
	vbox.PackStart(musicEndButton, false, false, 0)

	fileMvUpButton, err := gtk.ButtonNewWithLabel("Move Files Up")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	fileMvUpButton.SetMarginTop(5)
	fileMvUpButton.SetMarginBottom(5)
	fileMvUpButton.SetSizeRequest(150, -1)
	fileMvUpButton.Connect("clicked", func() {
		fileChooserDialog, err := gtk.FileChooserDialogNewWith2Buttons(
			"Select Folder",
			win,
			gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			"Cancel",
			gtk.RESPONSE_CANCEL,
			"Select",
			gtk.RESPONSE_ACCEPT,
		)
		if err != nil {
			log.Fatal("Unable to create file chooser dialog:", err)
		}
		if fileChooserDialog.Run() == gtk.RESPONSE_ACCEPT {
			folder := fileChooserDialog.GetFilename()
			fileChooserDialog.Destroy()
			fileMvUp(folder)
		} else {
			fileChooserDialog.Destroy()
		}
	})
	vbox.PackStart(fileMvUpButton, false, false, 0)

	ipListButton, err := gtk.ButtonNewWithLabel("List IP Addresses")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	ipListButton.SetMarginTop(5)
	ipListButton.SetMarginBottom(5)
	ipListButton.SetSizeRequest(150, -1)
	ipListButton.Connect("clicked", func() {
		cmd := exec.Command("sh", "-c", "ip a | awk '/^[0-9]+: / {iface=$2} /^[[:space:]]+inet / {print iface \" \" $2}'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to list IP addresses: %v\nOutput: %s", err, output)
		}
		formattedOutput := strings.TrimSpace(string(output))
		dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, formattedOutput)
		dialog.SetDefaultResponse(gtk.RESPONSE_OK)
		dialog.Run()
		dialog.Destroy()
	})
	vbox.PackStart(ipListButton, false, false, 0)

	fishEditButton, err := gtk.ButtonNewWithLabel("Edit Fish Config")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	fishEditButton.SetMarginTop(5)
	fishEditButton.SetMarginBottom(5)
	fishEditButton.SetSizeRequest(150, -1)
	fishEditButton.Connect("clicked", func() {
		fishEdit()
	})
	vbox.PackStart(fishEditButton, false, false, 0)

	win.Add(vbox)
	win.SetDefaultSize(400, 300)
	win.ShowAll()

	gtk.Main()
}

func mediaRestart() {
	cmd := exec.Command("systemctl", "--user", "restart", "wireplumber", "pipewire", "pipewire-pulse")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to restart media services: %v\nOutput: %s", err, output)
	}
}

var musicCmd *exec.Cmd

func musicStart() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	musicCmd = exec.Command("sh", "-c", "mpv --no-video --speed=1 --shuffle "+usr.HomeDir+"/Music/*")

	err = musicCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start music playback: %v", err)
	}

	log.Println("Music playback started successfully")
}

func musicEnd() {
	if musicCmd != nil && musicCmd.Process != nil {
		err := musicCmd.Process.Kill()
		if err != nil {
			log.Fatalf("Failed to stop music playback: %v", err)
		}
		log.Println("Music playback stopped successfully")
	} else {
		log.Println("No music playback process found")
	}
}

func fileMvUp(folder string) {
	cmd := exec.Command("sh", "-c", "find "+folder+" -mindepth 2 -type f -print -exec mv {} "+folder+" \\;")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to move files up: %v\nOutput: %s", err, output)
	}
	log.Printf("Files moved up successfully:\n%s", output)
}

func ipList() {
	cmd := exec.Command("sh", "-c", "ip -c a | awk '/^[0-9]+: / {print $2} /^[[:space:]]+inet / {print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to list IP addresses: %v\nOutput: %s", err, output)
	}
	fmt.Println(string(output))
}

func fishEdit() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	term := "gnome-terminal"
	termFlags := " -- bash -c"
	editor := "vim"
	filePath := usr.HomeDir + "/.config/fish/config.fish"

	cmd := exec.Command(term, termFlags, editor+" "+filePath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to edit fish config: %v", err)
	}

	fishReloadScript()
}

func fishReloadScript() {
	cmd := exec.Command("fish", "-c", "source /home/ghost/.config/fish/config.fish")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to reload fish config: %v", err)
	}
}
