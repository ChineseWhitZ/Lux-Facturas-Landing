package leads

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

var ErrNotFound = errors.New("lead not found")

type Repository struct {
	filePath string
	mu       sync.Mutex
	leads    []Lead
}

func NewRepository(filePath string) (*Repository, error) {
	repo := &Repository{
		filePath: filePath,
		leads:    []Lead{},
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *Repository) Create(input CreateLeadInput) (Lead, error) {
	if err := input.Validate(); err != nil {
		return Lead{}, err
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()

	lead := Lead{
		ID:                    newID(),
		CustomerName:          input.CustomerName,
		CompanyName:           input.CompanyName,
		DocumentNumber:        input.DocumentNumber,
		Phone:                 input.Phone,
		Email:                 input.Email,
		SelectedPlan:          input.SelectedPlan,
		CustomerType:          input.CustomerType,
		BusinessCategory:      input.BusinessCategory,
		ApproxProductQuantity: input.ApproxProductQuantity,
		City:                  input.City,
		Message:               input.Message,
		Status:                StatusNew,
		RegisteredAt:          time.Now().UTC(),
	}

	repo.leads = append(repo.leads, lead)

	if err := repo.saveLocked(); err != nil {
		return Lead{}, err
	}

	return lead, nil
}

func (repo *Repository) List() []Lead {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := slices.Clone(repo.leads)
	slices.SortFunc(result, func(a Lead, b Lead) int {
		return b.RegisteredAt.Compare(a.RegisteredAt)
	})

	return result
}

func (repo *Repository) FindByID(id string) (Lead, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, lead := range repo.leads {
		if lead.ID == id {
			return lead, nil
		}
	}

	return Lead{}, ErrNotFound
}

func (repo *Repository) UpdateStatus(id string, status Status) (Lead, error) {
	if !status.Valid() {
		return Lead{}, errors.New("status is invalid")
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()

	for index := range repo.leads {
		if repo.leads[index].ID == id {
			repo.leads[index].Status = status

			if err := repo.saveLocked(); err != nil {
				return Lead{}, err
			}

			return repo.leads[index], nil
		}
	}

	return Lead{}, ErrNotFound
}

func (repo *Repository) load() error {
	if repo.filePath == "" {
		return nil
	}

	content, err := os.ReadFile(repo.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return nil
	}

	return json.Unmarshal(content, &repo.leads)
}

func (repo *Repository) saveLocked() error {
	if repo.filePath == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(repo.filePath), 0755); err != nil {
		return err
	}

	content, err := json.MarshalIndent(repo.leads, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(repo.filePath, content, 0644)
}

func newID() string {
	return fmt.Sprintf("lead_%d", time.Now().UTC().UnixNano())
}
