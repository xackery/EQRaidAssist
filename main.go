package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	loadFile "github.com/Valorith/EQRaidAssist/loadFIle"
)

//Global Variables
var lastLoadedRaidFile string = ""
var raidScanner FileScanner
var logScanner FileScanner

func main() {
	loadedRaidFile, err := getNewestRaidFile()
	if err != nil {
		fmt.Println("lastLoadedRaidFile:", err)
	}
	loadedLogFile, err := getNewestLogFile()
	if err != nil {
		fmt.Println("loadedLogFile:", err)
	}
	lastLoadedRaidFile = loadedRaidFile

	raidScanner = FileScanner{loadedRaidFile, 1, false, 15} // Raid Dump Scan (Scan Type = 1)
	logScanner = FileScanner{loadedLogFile, 2, false, 15}   // Log Scan (Scan Type = 2)
	//Load Players in
	players, err := loadPlayers(lastLoadedRaidFile) // loads players from the newest Raid Dump file in the EQ directory
	if err != nil {
		fmt.Println("Player load failed: ", err)
		softExit()
	}

	playerCount := len(players)
	if playerCount > 0 {
		fmt.Printf("%d Players Succesfully loaded.\n", playerCount)
	} else {
		fmt.Println("Player Loading Failed!")
	}

	for _, player := range players {
		player.PrintPlayer()
	}

	//Get User user input
	getUserInput()
}

type Player struct {
	Name  string
	Level int
	Class string
	Group int
	Loot  []string
}

type Log struct {
	Directory string
	FileName  string
	FileSize  int64
	LastWrite time.Time
	LogData   []string
}

type FileScanner struct {
	fileLocation  string
	scanType      int
	enabled       bool
	ScanFrequency int //Frequency in seconds
}

func softExit() {
	var userInput string
	fmt.Println("Press enter to exit...")
	fmt.Scanln(&userInput)
	os.Exit(0)
}

func getUserInput() {
	var err error
	for {
		fmt.Printf("Commands:\nStart scanning raid file: 'start' or 'run'\nStop scanning raid file: 'stop'\nExit application: 'exit' or 'quit'\n")
		fmt.Println("-----------------")
		fmt.Println("Enter a command:")
		var userInput string
		fmt.Scanln(&userInput)
		switch userInput {
		case "start":
			if !raidScanner.IsRunning() {
				fileLocation := raidScanner.fileLocation
				if fileLocation != "" {
					raidScanner.SetScanFrequency(15)
					fmt.Println("Starting Up Raid Scanner...")
					err := raidScanner.Start()
					if err != nil {
						fmt.Println("error in raid scanner start:", err)
						continue
					}
				} else {
					err = fmt.Errorf("file not detected at location: %s", fileLocation)
					fmt.Println(err)
					continue
				}
			} else {
				err = fmt.Errorf("raidScanner is already running")
				fmt.Println(err)
				continue
			}
			if !logScanner.IsRunning() {
				fileLocation := logScanner.fileLocation
				if fileLocation != "" {
					logScanner.SetScanFrequency(15)
					fmt.Println("Starting Up Log Scanner...")
					err := logScanner.Start()
					if err != nil {
						fmt.Println("error in log scanner start:", err)
						continue
					}
				} else {
					err = fmt.Errorf("file not detected at location: %s", fileLocation)
					fmt.Println(err)
					continue
				}
			} else {
				err = fmt.Errorf("logScanner is already running")
				fmt.Println(err)
				continue
			}
		case "run":
			if !raidScanner.IsRunning() {
				fileLocation := raidScanner.fileLocation
				if fileLocation != "" {
					raidScanner.SetScanFrequency(15)
					fmt.Println("Starting Up Raid Scanner...")
					err := raidScanner.Start()
					if err != nil {
						fmt.Println("error in raid scanner start:", err)
						continue
					}
				} else {
					err = fmt.Errorf("file not detected at location: %s", fileLocation)
					fmt.Println(err)
					continue
				}
			} else {
				err = fmt.Errorf("raidScanner is already running")
				fmt.Println(err)
				continue
			}
			if !logScanner.IsRunning() {
				fileLocation := logScanner.fileLocation
				if fileLocation != "" {
					logScanner.SetScanFrequency(15)
					fmt.Println("Starting Up Log Scanner...")
					err := logScanner.Start()
					if err != nil {
						fmt.Println("error in log scanner start:", err)
						continue
					}
				} else {
					err = fmt.Errorf("file not detected at location: %s", fileLocation)
					fmt.Println(err)
					continue
				}
			} else {
				err = fmt.Errorf("logScanner is already running")
				fmt.Println(err)
				continue
			}
		case "stop":
			if raidScanner.IsRunning() {
				fmt.Println("[Status] Stopping Raid Scanner...")
				raidScanner.enabled = false
			} else {
				err = fmt.Errorf("the raid file scanner is not running")
				fmt.Println(err)
				continue
			}
			if logScanner.IsRunning() {
				fmt.Println("[Status] Stopping Log Scanner...")
				logScanner.enabled = false
			} else {
				err = fmt.Errorf("the log file scanner is not running")
				fmt.Println(err)
				continue
			}

		case "exit":
			fmt.Println("[Status] Exiting...")
			os.Exit(0)
		case "quit":
			fmt.Println("[Status] Exiting...")
			os.Exit(0)
		default:
			fmt.Println("[Status] Invalid command.")
		}
		if err != nil {
			fmt.Println(userInput, " failed: ", err)
			continue
		}
	}
}

func getNewestRaidFile() (string, error) {
	EQpath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("getwd: %w", err)
	}
	//logsFolder := EQpath + "\\Logs"
	//fmt.Println("Loading Players from: ", EQpath)
	raidDumpFileList := []string{}
	fileListIndex := 0
	filePathError := filepath.Walk(EQpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("filepath.Walk: %w", err)
		}
		//fmt.Printf("Scanning File: %s....\n", path)
		isDir := info.IsDir()
		itemPath := path
		if !isDir {
			if strings.Contains(itemPath, "RaidRoster") {
				raidDumpFileList = append(raidDumpFileList, itemPath)
				fileListIndex++
			}
		}
		return nil
	})

	if filePathError != nil {
		fmt.Println(filePathError)
		return "", fmt.Errorf("filepath.Walk: %w", filePathError)
	}

	newestIndex := len(raidDumpFileList) - 1
	return raidDumpFileList[newestIndex], nil
}

func getNewestLogFile() (string, error) {
	EQpath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("getwd: %w", err)
	}
	serverShortName, charName := "", ""

	//Get server short name
	userInput := ""
	fmt.Println("Enter the server short name:")
	fmt.Scanln(&userInput)
	serverShortName = userInput

	//Get character Name
	userInput = ""
	fmt.Println("Enter your characters first name:")
	fmt.Scanln(&userInput)
	charName = userInput

	//Set character log file path
	logsFolder := EQpath + "\\Logs\\eqlog_" + charName + "_" + serverShortName + ".txt"
	fmt.Println("Loading Log from: ", logsFolder)
	return logsFolder, err
}

func loadPlayers(fileLocation string) ([]Player, error) {
	fmt.Println("Newest Raid Dump File Detected: ", fileLocation)
	var dumpLines []string = loadFile.Load(fileLocation)
	var players []Player

	for _, line := range dumpLines {
		formattedLine := strings.Replace(line, "\t", ",", -1)
		groupNumber, _ := strconv.Atoi(formattedLine[0:strings.Index(formattedLine, ",")])
		formattedLine = formattedLine[strings.Index(formattedLine, ",")+1:]
		charName := formattedLine[0:strings.Index(formattedLine, ",")]
		formattedLine = formattedLine[strings.Index(formattedLine, ",")+1:]
		charLevel, _ := strconv.Atoi(formattedLine[0:strings.Index(formattedLine, ",")])
		formattedLine = formattedLine[strings.Index(formattedLine, ",")+1:]
		charClass := formattedLine[0:strings.Index(formattedLine, ",")]
		fmt.Println("Player Detected: ", charName)
		players = append(players, Player{charName, charLevel, charClass, groupNumber, []string{}})
	}

	return players, nil
}

func loadLog(fileLocation string) (Log, error) {
	fmt.Println("Newest Raid Dump File Detected: ", fileLocation)
	fileSize, err := loadFile.GetFileSize(fileLocation)
	if err != nil {
		return Log{}, fmt.Errorf("getFileSize: %w", err)
	}
	logLines := loadFile.Load(fileLocation)
	log := Log{}
	log.Directory = fileLocation
	log.FileSize = fileSize

	if len(logLines) == 0 {
		err = fmt.Errorf("failed to load log at location: %s", fileLocation)
		return Log{}, err
	}
	log.ClearData()
	for _, line := range logLines {
		log.LogData = append(log.LogData, line)
		fmt.Println(line)
	}

	return log, nil
}

func (player *Player) PrintPlayer() {

	fmt.Println("Char Name: ", player.Name)
	fmt.Println("Char Level: ", player.Level)
	fmt.Println("Char Class: ", player.Class)
	fmt.Println("Group Number: ", player.Group)
	fmt.Println("Loot:")
	for _, lootItem := range player.Loot {
		fmt.Println("\t", lootItem)
	}
	fmt.Println("-----------------")

}

func (player *Player) AddLoot(lootItem string) {
	if lootItem != "" {
		player.Loot = append(player.Loot, lootItem)
	} else {
		fmt.Println("Error adding loot item: ", lootItem)
	}
}

func (scanner *FileScanner) scan() error {
	// Ensure Scanner fileLocation is set
	if scanner.fileLocation == "" {
		return fmt.Errorf("fileLocation is not set")
	}
	if !scanner.IsRunning() {
		scanner.enabled = true // Ensure scanner is enabled
	}
	for {
		if scanner.IsRunning() {
			switch scanner.GetType() {
			case 1:
				//Scan Raid Dump file
				fmt.Println("Scanning File:", scanner.fileLocation, "with a frequency of", scanner.ScanFrequency, "seconds")
				loadedRaidFile, err := getNewestRaidFile()
				if err != nil {
					return err
				}

				if loadedRaidFile != lastLoadedRaidFile {
					//Load Players in
					players, err := loadPlayers(loadedRaidFile) // loads players from the newest Raid Dump file in the EQ directory
					if err != nil {
						fmt.Println("Player load failed: ", err)
					}
					playerCount := len(players)
					if playerCount > 0 {
						fmt.Printf("%d Players Succesfully loaded.\n", playerCount)
					} else {
						fmt.Println("Player Loading Failed!")
					}

					for _, player := range players {
						player.PrintPlayer()
					}
					lastLoadedRaidFile = loadedRaidFile
				} else {
					fmt.Println("No new raid file detected.")
				}
				time.Sleep(time.Second * time.Duration(scanner.ScanFrequency))
			case 2:
				//Scan character log file
				fmt.Println("Scanning File:", scanner.fileLocation, "with a frequency of", scanner.ScanFrequency, "seconds")
				logFileDir, err := getNewestLogFile()
				if err != nil {
					fmt.Println("scan error: getNewestLogFile: %w", err)
					continue
				}
				//Load log data in
				log, err := loadLog(logFileDir) // loads players from the newest Raid Dump file in the EQ directory
				if err != nil {
					fmt.Println("scan error: loadLog: %w", err)
					continue
				}
				logLength := len(log.LogData)
				if logLength > 0 {
					fmt.Printf("%d Log entries succesfully loaded.\n", logLength)
				} else {
					fmt.Println("log loading failed. length = ", logLength)
					continue
				}

				for _, logEntry := range log.LogData {
					fmt.Println(logEntry)
				}

				time.Sleep(time.Second * time.Duration(scanner.ScanFrequency))
			}
		}
	}
}

func (scanner *FileScanner) IsRunning() bool {
	return scanner.enabled
}

func (scanner *FileScanner) Start() error {
	if !scanner.IsRunning() {
		if scanner.fileLocation != "" {
			scanChannel := make(chan error, 2)
			go scanner.scan()
			err := <-scanChannel
			if err != nil {
				return fmt.Errorf("scan error: %w", err)
			}
		} else {
			return fmt.Errorf("file location has not been set for the scanner")
		}
	} else {
		return fmt.Errorf("scanner is already running")
	}
	return nil
}

func (scanner *FileScanner) SetType(scanType int) {
	scanner.scanType = scanType
}

func (scanner *FileScanner) GetType() int {
	return scanner.scanType
}

func (scanner *FileScanner) SetScanFrequency(frequency int) {
	scanner.ScanFrequency = frequency
}

func (log *Log) ClearData() {
	log.LogData = []string{}
}
