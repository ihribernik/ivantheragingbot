package soundboard

type Service struct {
	resolver AssetResolver
	player   AudioPlayer
}

func New(resolver AssetResolver, player AudioPlayer) *Service {
	return &Service{
		resolver: resolver,
		player:   player,
	}
}

func (s *Service) Play(key string) error {
	path, err := s.resolver.Resolve(key)
	if err != nil {
		return err
	}

	return s.player.Play(path)
}
