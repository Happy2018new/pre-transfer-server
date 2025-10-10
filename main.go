package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Happy2018new/pre-transfer-server/minecraft"
	"github.com/Happy2018new/pre-transfer-server/minecraft/protocol"
	"github.com/Happy2018new/pre-transfer-server/minecraft/protocol/packet"
	"github.com/Happy2018new/pre-transfer-server/minecraft/resource"
)

// Config ..
type Config struct {
	Local  ServerAddress `json:"local"`
	Remote ServerAddress `json:"remote"`
}

// ServerAddress ..
type ServerAddress struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

func main() {
	var cfg Config
	var packs []*resource.Pack

	entries, err := os.ReadDir("packs")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileBytes, err := os.ReadFile(filepath.Join("packs", entry.Name()))
		if err != nil {
			panic(err)
		}

		buf := bytes.NewBuffer(fileBytes)
		pack, err := resource.Read(buf)
		if err != nil {
			panic(err)
		}

		packs = append(packs, pack)
	}

	fileBytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileBytes, &cfg)
	if err != nil {
		panic(err)
	}

	config := minecraft.ListenConfig{
		AuthenticationDisabled: true,
		ResourcePacks:          packs,
	}
	listener, err := config.Listen("raknet", fmt.Sprintf("%s:%d", cfg.Local.Address, cfg.Local.Port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		netConn, err := listener.Accept()
		if err != nil {
			continue
		}
		conn := netConn.(*minecraft.Conn)

		err = conn.StartGame(minecraft.GameData{
			WorldName:           "World",
			WorldSeed:           0,
			Difficulty:          0, // Peaceful
			EntityUniqueID:      0,
			EntityRuntimeID:     0,
			PlayerGameMode:      1, // Creative
			PersonaDisabled:     false,
			CustomSkinsDisabled: false,
			BaseGameVersion:     "*",
			PlayerPosition:      [3]float32{0.5, 1.5, 0.5},
			Pitch:               0,
			Yaw:                 0,
			Dimension:           0, // Overworld
			WorldSpawn:          [3]int32{0, 0, 0},
			EditorWorldType:     packet.EditorWorldTypeNotEditor,
			CreatedInEditor:     false,
			WorldGameMode:       1, // Creative
			Hardcore:            false,
			GameRules: []protocol.GameRule{
				{
					Name:  "doDayLightCycle",
					Value: false,
				},
			},
			Time:                     0,
			ServerBlockStateChecksum: 0,
			CustomBlocks:             nil,
			Items:                    nil,
			PlayerMovementSettings: protocol.PlayerMovementSettings{
				RewindHistorySize:                0,
				ServerAuthoritativeBlockBreaking: false,
			},
			ServerAuthoritativeInventory: true,
			Experiments:                  nil,
			PlayerPermissions:            1, // Member
			ChunkRadius:                  4,
			ClientSideGeneration:         false,
			ChatRestrictionLevel:         packet.ChatRestrictionLevelDisabled,
			DisablePlayerInteractions:    true,
			UseBlockNetworkIDHashes:      false,
		})
		if err != nil {
			_ = conn.Close()
		}

		_ = conn.WritePacket(&packet.Transfer{
			Address: cfg.Remote.Address,
			Port:    cfg.Remote.Port,
		})
		_ = conn.Close()
	}
}
