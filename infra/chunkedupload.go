package infra

type ChunkUploadCallback func(chunkId int, data []byte) error
