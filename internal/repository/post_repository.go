package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostRepository struct {
	DB *sql.DB
}

type RepositoryPost interface {
	CreatePost(ctx context.Context, request *model.CreatePostRequest) (*model.PostResponse, error)
	FindPostById(ctx context.Context, id string) (*model.PostResponse, error)
	PostList(ctx context.Context, request *model.PostListRequest) ([]model.PostListResponse, model.MetaDataResponse, error)
	CreatePostComment(ctx context.Context, request *model.CreatePostCommentRequest) (*model.PostCommentResponse, error)
}

func NewPostRepository(db *sql.DB) RepositoryPost {
	return &PostRepository{
		DB: db,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, request *model.CreatePostRequest) (*model.PostResponse, error) {

	var post model.PostResponse

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `INSERT INTO posts (user_id, content, tags, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, user_id, content, tags`

	err := r.DB.QueryRowContext(context, query, request.UserId, request.Content, request.Tags, time.Now(), time.Now()).Scan(&post.Id, &post.UserId, &post.Content, &post.Tags)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) FindPostById(ctx context.Context, id string) (*model.PostResponse, error) {
	var post model.PostResponse

	validate := uuid.Validate(id)
	if validate != nil {
		return nil, model.ErrResNotFound.Error
	}

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `SELECT id, user_id, content, tags, created_at, updated_at FROM posts WHERE id = $1`

	err := r.DB.QueryRowContext(context, query, id).Scan(&post.Id, &post.UserId, &post.Content, &post.Tags, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrResNotFound.Error
		}
		return nil, model.ErrInternalDatabase
	}

	return &post, nil
}

func (r *PostRepository) CreatePostComment(ctx context.Context, request *model.CreatePostCommentRequest) (*model.PostCommentResponse, error) {

	var post model.PostCommentResponse

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `INSERT INTO post_comments (post_id, user_id, comment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, post_id, user_id, comment`

	err := r.DB.QueryRowContext(context, query, request.PostId, request.UserId, request.Comment, time.Now(), time.Now()).Scan(&post.Id, &post.PostId, &post.UserId, &post.Comment)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) PostList(ctx context.Context, request *model.PostListRequest) ([]model.PostListResponse, model.MetaDataResponse, error) {

	var posts []model.PostListResponse
	var metaData model.MetaDataResponse

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// queryCondition := fmt.Sprintf("WHERE posts.userd_id <> '%s' ", request.UserId)
	queryCondition := fmt.Sprintf("")
	if request.Search != "" {
		if queryCondition == "" {
			queryCondition += " WHERE"
		} else {
			queryCondition += " AND"
		}
		queryCondition += " posts.content LIKE '%" + request.Search + "%'"
	}

	if len(request.SearchTag) > 0 {
		jsonTag, err := json.Marshal([]string(request.SearchTag))
		if err == nil {
			replacer := strings.NewReplacer("[", "{", "]", "}")
			stringTag := replacer.Replace(string(jsonTag))

			if queryCondition == "" {
				queryCondition += " WHERE "
			} else {
				queryCondition += " AND "
			}

			queryCondition += fmt.Sprintf(" posts.tags && '%s'", stringTag)
		}
	}

	queryGet := fmt.Sprintf(`SELECT DISTINCT
	posts."id",
	posts.user_id,
	posts."content",
	posts.tags,
	posts.created_at,
	users.ID AS post_creator_user_id,
	users."name" AS post_creator_user_name,
	users.total_friend AS post_creator_total_friend,
	users.image_url AS post_creator_image_url,
	ARRAY (
	SELECT
		(
			post_comments."id" || ',' || post_comments.post_id || ',' || post_comments."user_id" || ',' || post_comments.COMMENT || ',' || post_comments.created_at || ',' || users."id"  || ',' || users.NAME || ',' || users.image_url  || ',' || users.total_friend || ',' || users.created_at 
		) 
	FROM
		post_comments
		JOIN users ON post_comments.user_id = users.ID 
	WHERE
		posts.ID = post_comments.post_id 
	) AS comments_post 
FROM
	posts
	JOIN users ON posts.user_id = users.
	ID LEFT JOIN post_comments ON post_comments.post_id = posts."id" 
    %s
ORDER BY
	posts.created_at DESC 
	LIMIT %d OFFSET %d;`, queryCondition, request.Limit, request.Offset)

	rows, err := r.DB.QueryContext(context, queryGet)
	if err != nil {
		return nil, model.MetaDataResponse{}, err
	}

	totalRows := 0
	for rows.Next() {

		var post model.PostListResponse
		var createdAt time.Time
		postCommentString := make([]string, 0)
		err := rows.Scan(
			&post.PostId,
			&post.Post.UserId,
			&post.Post.Content,
			&post.Post.Tags,
			&createdAt,
			&post.Creator.UserId,
			&post.Creator.Name,
			&post.Creator.FriendCount,
			&post.Creator.ImageUrl,
			pq.Array(&postCommentString),
		)
		if err != nil {
			return nil, model.MetaDataResponse{}, err
		}

		post.Post.CreatedAt = createdAt

		postComment := make([]model.PostCommentUserResponse, 0)
		for i := 0; i < (len(postCommentString)); i++ {
			cArray := strings.Split(postCommentString[i], ",")

			var comment model.PostCommentUserResponse

			comment.Id = cArray[0]
			comment.PostId = cArray[1]
			comment.UserId = cArray[2]
			comment.Comment = cArray[3]

			comment.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", cArray[4])

			comment.Creator.UserId = cArray[5]
			comment.Creator.Name = cArray[6]
			comment.Creator.ImageUrl = cArray[7]
			friendCount, _ := strconv.Atoi(cArray[8])
			comment.Creator.FriendCount = friendCount
			comment.Creator.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", cArray[9])

			postComment = append(postComment, comment)
		}

		post.Comments = postComment

		posts = append(posts, post)

		totalRows++
	}

	rows.Close()

	return posts, metaData, nil
}
