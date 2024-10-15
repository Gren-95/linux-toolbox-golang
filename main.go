package main

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Authentication information
const (
	vpnUser    = "" // Replace with actual VPN user
	vpnAddress = "" // Replace with actual VPN address
	vpnGroup   = "" // Replace with actual VPN group
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

	var rpgmCheatsCmd = &cobra.Command{
		Use:   "rpgmCheats",
		Short: "Clone RPGM Cheat Menu and run the patcher",
		Run: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			if path == "" {
				var err error
				path, err = os.Getwd()
				if err != nil {
					log.Fatalf("Failed to get current directory: %v", err)
				}
			}
			rpgmCheats(path)
		},
	}
	rpgmCheatsCmd.Flags().StringP("path", "p", "", "Path to the directory")

	var ooklaSpeedTestCmd = &cobra.Command{
		Use:   "ooklaSpeedTest",
		Short: "Run Ookla Speed Test",
		Run: func(cmd *cobra.Command, args []string) {
			ooklaSpeedTest()
		},
	}

	var timeshiftBackupCmd = &cobra.Command{
		Use:   "timeshiftBackup",
		Short: "Create a Timeshift backup",
		Run: func(cmd *cobra.Command, args []string) {
			timeshiftBackup()
		},
	}
	var jeldwenVpnConnectCmd = &cobra.Command{
		Use:   "jeldwenVpnConnect",
		Short: "Connect to Jeld-Wen VPN",
		Run: func(cmd *cobra.Command, args []string) {
			jeldwenVpnConnect()
		},
	}

	rootCmd.AddCommand(jeldwenVpnConnectCmd)
	rootCmd.AddCommand(timeshiftBackupCmd)
	rootCmd.AddCommand(ooklaSpeedTestCmd)
	rootCmd.AddCommand(rpgmCheatsCmd)
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

	rpgmCheatsButton, err := gtk.ButtonNewWithLabel("Run RPGM Cheats")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	rpgmCheatsButton.SetMarginTop(5)
	rpgmCheatsButton.SetMarginBottom(5)
	rpgmCheatsButton.SetSizeRequest(150, -1)
	rpgmCheatsButton.Connect("clicked", func() {
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
			rpgmCheats(folder)
		} else {
			fileChooserDialog.Destroy()
		}
	})
	vbox.PackStart(rpgmCheatsButton, false, false, 0)

	ooklaSpeedTestButton, err := gtk.ButtonNewWithLabel("Run Speed Test")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	ooklaSpeedTestButton.SetMarginTop(5)
	ooklaSpeedTestButton.SetMarginBottom(5)
	ooklaSpeedTestButton.SetSizeRequest(150, -1)
	ooklaSpeedTestButton.Connect("clicked", func() {
		ooklaSpeedTest()
	})
	vbox.PackStart(ooklaSpeedTestButton, false, false, 0)

	timeshiftBackupButton, err := gtk.ButtonNewWithLabel("Create Timeshift Backup")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	timeshiftBackupButton.SetMarginTop(5)
	timeshiftBackupButton.SetMarginBottom(5)
	timeshiftBackupButton.SetSizeRequest(200, -1)
	timeshiftBackupButton.Connect("clicked", func() {
		go func() {
			timeshiftBackup()
			glib.IdleAdd(func() {
				dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "Backup complete")
				dialog.SetDefaultResponse(gtk.RESPONSE_OK)
				dialog.Run()
				dialog.Destroy()
			})
		}()
	})
	vbox.PackStart(timeshiftBackupButton, false, false, 0)

	jeldwenVpnConnectButton, err := gtk.ButtonNewWithLabel("Connect to Jeld-Wen VPN")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	jeldwenVpnConnectButton.SetMarginTop(5)
	jeldwenVpnConnectButton.SetMarginBottom(5)
	jeldwenVpnConnectButton.SetSizeRequest(200, -1)
	jeldwenVpnConnectButton.Connect("clicked", func() {
		jeldwenVpnConnect()
	})
	vbox.PackStart(jeldwenVpnConnectButton, false, false, 0)

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
		log.Fatalf("Failed to get current user: %v", err)
	}

	musicCmd = exec.Command("mpv", "--no-video", "--speed=1", "--shuffle", usr.HomeDir+"/Music/*")
	err = musicCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start music playback: %v", err)
	}

	log.Printf("Music playback started successfully with command: %v", musicCmd.Args)
}

func musicEnd() {
	if musicCmd != nil && musicCmd.Process != nil {
		err := musicCmd.Process.Signal(os.Interrupt)
		if err != nil {
			log.Fatalf("Failed to stop music playback: %v", err)
		}
		log.Println("Music playback stopped successfully")
	} else {
		log.Println("No music playback process found")
	}
}

func fileMvUp(folder string) {
	cmd := exec.Command("sh", "-c", "find \""+folder+"\" -mindepth 2 -type f -print -exec mv {} \""+folder+"\" \\;")
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
		log.Fatalf("Failed to get current user: %v", err)
	}

	term := "gnome-terminal"
	termFlags := []string{"--", "bash", "-c"}
	editor := "vim \"" + usr.HomeDir + "/.config/fish/config.fish\""
	cmd := exec.Command(term, append(termFlags, editor)...)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run command '%s %s': %v", term, strings.Join(termFlags, " "), err)
	}

	fishReloadScript()
}

func fishReloadScript() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	cmd := exec.Command("fish", "-c", "source \""+usr.HomeDir+"/.config/fish/config.fish\"")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to reload fish config: %v", err)
	}
}

func rpgmCheats(path string) {
	// Check if wine is installed
	_, err := exec.LookPath("wine")
	if err != nil {
		log.Fatalf("wine is not installed or not found in PATH: %v", err)
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf(`
if [ ! -d "RPGM_Cheat_Menu" ]; then
    git clone https://github.com/Gren-95/RPGM_Cheat_Menu.git
fi
cp -r RPGM_Cheat_Menu/MVPluginPatcher.exe RPGM_Cheat_Menu/plugins_patch.txt RPGM_Cheat_Menu/www/ "%s"
`, path))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run RPGM Cheats: %v\nOutput: %s", err, output)
	}
	log.Printf("RPGM Cheats executed successfully:\nDo not forget to actually run the MVPluginPatcher.exe file!\n%s", output)
}

func ooklaSpeedTest() {

	term := "gnome-terminal"
	termFlags := []string{"--", "bash", "-c"}
	command := "speedtest-cli && sleep 5"
	cmd := exec.Command(term, append(termFlags, command)...)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run speedtest-cli: %v", err)
	}
}

func timeshiftBackup() {
	cmd := exec.Command("sudo", "timeshift", "--create")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to create timeshift backup: %v", err)
	}
}

func jeldwenVpnConnect() {
	cmd := exec.Command("gnome-terminal", "--", "sudo", "openconnect", "--user="+vpnUser, vpnAddress, "--authgroup="+vpnGroup)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to connect to VPN: %v", err)
	}
}
