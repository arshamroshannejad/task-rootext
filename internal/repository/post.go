package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arshamroshannejad/task-rootext/internal/domain"
	"github.com/arshamroshannejad/task-rootext/internal/entities"
	"github.com/arshamroshannejad/task-rootext/internal/helpers"
	"github.com/arshamroshannejad/task-rootext/internal/model"
	"time"
)

type postRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) domain.PostRepository {
	return &postRepositoryImpl{
		db: db,
	}
}

func (p *postRepositoryImpl) GetAll(filter *helpers.PaginateFilter) (*[]model.Post, helpers.Metadata, error) {
	query := fmt.Sprintf(
		`
			SELECT 
				COUNT(*) OVER() AS total_records,
				p.id, 
				p.title, 
				p.text, 
				p.created_at, 
				p.updated_at, 
				p.user_id, 
				COALESCE(SUM(v.vote), 0) AS vote_count
			FROM 
				posts p
			LEFT JOIN 
				votes v ON p.id = v.post_id
			GROUP BY 
				p.id
			ORDER BY 
				%s %s
			LIMIT
				$1 
			OFFSET 
				$2;
        `,
		filter.SortValue(),
		filter.SortDirection(),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, query, filter.Limit(), filter.OffSet())
	if err != nil {
		return nil, helpers.Metadata{}, err
	}
	defer rows.Close()
	return collectPostRows(rows, filter.Limit(), filter.OffSet())
}

func (p *postRepositoryImpl) GetByID(postID string) (*model.Post, error) {
	query := `
			SELECT 
				p.id,
				p.title,
				p.text,
				p.created_at,
				p.updated_at,
				p.user_id,
				COALESCE(SUM(v.vote), 0) as vote_count
			FROM posts p
			LEFT JOIN 
			    votes v ON p.id = v.post_id
			WHERE 
			    p.id = $1
			GROUP BY 
			    p.id
        `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	row := p.db.QueryRowContext(ctx, query, postID)
	return collectPostRow(row)
}

func (p *postRepositoryImpl) GetByTitle(title string) (*model.Post, error) {
	query := `
                SELECT 
                    p.id, p.title, p.text, p.created_at, p.updated_at, p.user_id, COALESCE(SUM(v.vote), 0) as vote_count
                FROM posts p
                LEFT JOIN votes v ON p.id = v.post_id
                WHERE p.title = $1
                GROUP BY p.id
        `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	row := p.db.QueryRowContext(ctx, query, title)
	return collectPostRow(row)
}

func (p *postRepositoryImpl) Create(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error) {
	query := `
                INSERT INTO posts (title, text, user_id) 
                VALUES ($1, $2, $3) 
                RETURNING id, title, text, created_at, updated_at, user_id, 0 as vote_count
        `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	args := []any{post.Title, post.Text, userID}
	row := p.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return collectPostRow(row)
}

func (p *postRepositoryImpl) Update(post *entities.PostCreateUpdateRequest, postID string) (*model.Post, error) {
	query := `
                UPDATE posts 
                SET title = $1, text = $2, updated_at = CURRENT_TIMESTAMP 
                WHERE id = $3 
                RETURNING id, title, text, created_at, updated_at, user_id, 0 as vote_count
        `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	args := []any{post.Title, post.Text, postID}
	row := p.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return collectPostRow(row)
}

func (p *postRepositoryImpl) Delete(postID string) error {
	query := "DELETE FROM posts WHERE id = $1"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := p.db.ExecContext(ctx, query, postID)
	return err
}

func (p *postRepositoryImpl) AddVote(postID, userID, vote string) error {
	query := "INSERT INTO votes (user_id, post_id, vote) VALUES ($1, $2, $3) ON CONFLICT (user_id, post_id) DO UPDATE SET vote = $4"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	args := []any{postID, userID, vote, vote}
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *postRepositoryImpl) RemoveVote(postID, userID string) error {
	query := "DELETE FROM votes WHERE user_id = $1 AND post_id = $2"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	args := []any{postID, userID}
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func collectPostRows(rows *sql.Rows, limit, offset int) (*[]model.Post, helpers.Metadata, error) {
	var posts []model.Post
	var totalRecords int
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&totalRecords,
			&post.ID,
			&post.Title,
			&post.Text,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UserID,
			&post.VoteCount,
		)
		if err != nil {
			return nil, helpers.Metadata{}, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, helpers.Metadata{}, err
	}
	metadata := helpers.CalculateMetadata(totalRecords, offset, limit)
	return &posts, metadata, nil
}

func collectPostRow(row *sql.Row) (*model.Post, error) {
	var post model.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Text,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.VoteCount,
	)
	if err != nil {
		return nil, err
	}
	return &post, err
}
