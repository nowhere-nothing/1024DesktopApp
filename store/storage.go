package store

type Storage interface {
	Save(pc *PostContent, pi *PostImage) error
	SaveFailed(pc *PostContent, failed string) error
	MkdirAll(pc *PostContent) error
}
