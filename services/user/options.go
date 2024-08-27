package user

import "github.com/spf13/viper"

func WithConfig(config *viper.Viper) func(s *Service) error {
	return func(s *Service) error {
		s.config = config
		return nil
	}
}
