package soundboard

type AudioPlayer interface {
	Play(path string) error
}

type AssetResolver interface {
	Resolve(key string) (string, error)
}
