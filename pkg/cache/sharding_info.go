package cache

type ShardingType int

const (
	ShardingTypeJWT ShardingType = iota
	ShardingTypeInteger
)

type ShardingInfo struct {
	Shards   int
	HashFunc func(key string) int
}

func getShardingInfo(shardingType ShardingType) ShardingInfo {
	switch shardingType {
	case ShardingTypeInteger:
		return ShardingInfo{
			Shards: 10,
			HashFunc: func(key string) int {
				char := int(key[len(key)-1])
				return char - 0x30
			},
		}
	case ShardingTypeJWT:
		return ShardingInfo{
			Shards: 36,
			HashFunc: func(key string) int {
				char := int(key[len(key)-1])
				if char >= 0x30 && char <= 0x39 {
					return char - 0x30
				}
				if char >= 0x41 && char <= 0x5A {
					return char - 0x41
				}
				return char - 0x61
			},
		}
	}

	return ShardingInfo{
		Shards: 1,
		HashFunc: func(key string) int {
			return 0
		},
	}
}
