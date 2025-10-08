package main

import (
	"encoding/json"
	"os"

	"github.com/Happy2018new/pre-transfer-server/minecraft"
	"github.com/Happy2018new/pre-transfer-server/minecraft/protocol"
	"github.com/Happy2018new/pre-transfer-server/minecraft/protocol/packet"
)

// RemoteServerAddr ..
type RemoteServerAddr struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

func main() {
	var serverAddress RemoteServerAddr

	fileBytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileBytes, &serverAddress)
	if err != nil {
		panic(err)
	}

	config := minecraft.ListenConfig{
		AuthenticationDisabled: true,
	}
	listener, err := config.Listen("raknet", "127.0.0.1:2025")
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
			Address: serverAddress.Address,
			Port:    serverAddress.Port,
		})
		_ = conn.Close()
	}
}
