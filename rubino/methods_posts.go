package rubino

import "context"

// ===== Posts (likes, comments, views, fetch) =========================

// LikePostAction likes a post.
func (s *Rubino) LikePostAction(ctx context.Context, profileID, postID, postProfileID, actionType string) (*Event, error) {
	return s.request(ctx, "likePostAction", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"action_type":     actionType,
	})
}

// Like is a convenience wrapper for LikePostAction with action="Like".
func (s *Rubino) Like(ctx context.Context, postID, postProfileID string) (*Event, error) {
	return s.LikePostAction(ctx, "", postID, postProfileID, "Like")
}

// Unlike is a convenience wrapper for LikePostAction with action="Unlike".
func (s *Rubino) Unlike(ctx context.Context, postID, postProfileID string) (*Event, error) {
	return s.LikePostAction(ctx, "", postID, postProfileID, "Unlike")
}

// GetComments returns the comments on a post.
func (s *Rubino) GetComments(ctx context.Context, profileID, postID, postProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getComments", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"limit":           limit,
		"sort":            ifEmpty(sort, "latest"),
		"equal":           equal,
	})
}

// AddComment comments on a post.
func (s *Rubino) AddComment(ctx context.Context, profileID, content, postID, postProfileID string) (*Event, error) {
	return s.request(ctx, "addComment", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"content":         content,
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"rnd":             rnd(),
	})
}

// AddReplyComment replies to a comment on a post.
func (s *Rubino) AddReplyComment(ctx context.Context, profileID, content, commentID, postID string) (*Event, error) {
	return s.request(ctx, "addReplyComment", map[string]any{
		"profile_id": s.profileOr(profileID),
		"content":    content,
		"comment_id": commentID,
		"post_id":    postID,
		"rnd":        rnd(),
		"track_id":   "Explore:two_tower",
	})
}

// LikeComment likes a comment.
func (s *Rubino) LikeComment(ctx context.Context, profileID, commentID, postID string) (*Event, error) {
	return s.request(ctx, "likeCommentAction", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"comment_id":  commentID,
		"post_id":     postID,
		"action_type": "Like",
	})
}

// UnlikeComment unlikes a comment.
func (s *Rubino) UnlikeComment(ctx context.Context, profileID, commentID, postID string) (*Event, error) {
	return s.request(ctx, "likeCommentAction", map[string]any{
		"profile_id":  s.profileOr(profileID),
		"comment_id":  commentID,
		"post_id":     postID,
		"action_type": "Unlike",
	})
}

// GetCommentReplies returns the replies to a comment.
func (s *Rubino) GetCommentReplies(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*Event, error) {
	return s.request(ctx, "getCommentReplies", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"comment_id": commentID,
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
	})
}

// DeleteComment deletes a comment. The removeRecord API expects a
// `record_id` + `model` pair (model="Comment"); the post_id is also
// required so the server can locate the comment.
func (s *Rubino) DeleteComment(ctx context.Context, profileID, commentID, postID string) (*Event, error) {
	return s.request(ctx, "removeRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"record_id":  commentID,
		"model":      "Comment",
		"post_id":    postID,
	})
}

// GetPostLikes returns the users who liked a post.
func (s *Rubino) GetPostLikes(ctx context.Context, profileID, postID, postProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getPostLikes", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"limit":           limit,
		"sort":            ifEmpty(sort, "latest"),
		"equal":           equal,
	})
}

// GetLikedCommentProfiles returns the profiles who liked a comment.
func (s *Rubino) GetLikedCommentProfiles(ctx context.Context, profileID, postID, commentID string, limit int, sort string) (*Event, error) {
	return s.request(ctx, "getLikedCommentProfiles", map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_id":    postID,
		"comment_id": commentID,
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
	})
}

// AddPostViewCount registers a view on a post.
func (s *Rubino) AddPostViewCount(ctx context.Context, postID, postProfileID string) (*Event, error) {
	return s.request(ctx, "addPostViewCount", map[string]any{
		"post_id":         postID,
		"post_profile_id": postProfileID,
	})
}

// AddPostViewTime registers a view duration on a post.
func (s *Rubino) AddPostViewTime(ctx context.Context, postID, postProfileID string, duration int) (*Event, error) {
	return s.request(ctx, "addPostViewTime", map[string]any{
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"duration":        duration,
		"profile_id":      s.profile,
	})
}

// GetRecentFollowingPosts returns the feed of recent posts from
// followings.
func (s *Rubino) GetRecentFollowingPosts(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getRecentFollowingPosts", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
		"equal":      equal,
	})
}

// GetProfilePosts returns the posts of a specific target profile.
func (s *Rubino) GetProfilePosts(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getProfilePosts", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
	})
}

// GetBookmarkedPosts returns the user's saved posts.
func (s *Rubino) GetBookmarkedPosts(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getBookmarkedPosts", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
		"equal":      equal,
	})
}

// SavePost bookmarks a post.
func (s *Rubino) SavePost(ctx context.Context, profileID, postID, postProfileID string) (*Event, error) {
	return s.request(ctx, "postBookmarkAction", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"action_type":     "Bookmark",
	})
}

// UnsavePost removes a bookmark.
func (s *Rubino) UnsavePost(ctx context.Context, profileID, postID, postProfileID string) (*Event, error) {
	return s.request(ctx, "postBookmarkAction", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"action_type":     "Unmark",
	})
}

// GetExplorePosts returns the explore feed.
func (s *Rubino) GetExplorePosts(ctx context.Context, profileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getExplorePosts", map[string]any{
		"profile_id": s.profileOr(profileID),
		"limit":      limit,
		"sort":       ifEmpty(sort, "latest"),
		"equal":      equal,
	})
}

// GetRelatedExplorePost returns related explore posts.
func (s *Rubino) GetRelatedExplorePost(ctx context.Context, profileID, postID, postProfileID string, limit int) (*Event, error) {
	return s.request(ctx, "getRelatedExplorePost", map[string]any{
		"profile_id":      s.profileOr(profileID),
		"post_id":         postID,
		"post_profile_id": postProfileID,
		"limit":           limit,
	})
}

// GetExplorePostTopics returns the explore post topics.
func (s *Rubino) GetExplorePostTopics(ctx context.Context, profileID string) (*Event, error) {
	return s.request(ctx, "getExplorePostTopics", map[string]any{
		"profile_id": s.profileOr(profileID),
	})
}

// GetTaggedPosts returns the tagged posts of a target profile.
func (s *Rubino) GetTaggedPosts(ctx context.Context, profileID, targetProfileID string, limit int, sort, equal string) (*Event, error) {
	return s.request(ctx, "getTaggedPosts", map[string]any{
		"profile_id":        s.profileOr(profileID),
		"target_profile_id": targetProfileID,
		"limit":             limit,
		"sort":              ifEmpty(sort, "latest"),
		"equal":             equal,
	})
}

// GetPostByShareLink resolves a `rubika.ir/post/...` share link.
func (s *Rubino) GetPostByShareLink(ctx context.Context, profileID, shareLink string) (*Event, error) {
	return s.request(ctx, "getPostByShareLink", map[string]any{
		"profile_id": s.profileOr(profileID),
		"url_post":   shareLink,
	})
}

// DeletePost deletes a post. The removeRecord API expects a
// `record_id` + `model` pair (model="Post").
func (s *Rubino) DeletePost(ctx context.Context, profileID, postID string) (*Event, error) {
	return s.request(ctx, "removeRecord", map[string]any{
		"profile_id": s.profileOr(profileID),
		"record_id":  postID,
		"model":      "Post",
	})
}

// AddPost publishes a new post. `postType` is "Picture" or "Video".
// The file must already be uploaded (see RequestUploadFile + UploadFile).
func (s *Rubino) AddPost(ctx context.Context, profileID, postType, fileID, caption, fileName string, width, height int) (*Event, error) {
	in := map[string]any{
		"profile_id": s.profileOr(profileID),
		"post_type":  postType,
		"file_id":    fileID,
		"rnd":        rnd(),
	}
	if caption != "" {
		in["caption"] = caption
	}
	if fileName != "" {
		in["file_name"] = fileName
	}
	if width > 0 {
		in["width"] = width
	}
	if height > 0 {
		in["height"] = height
	}
	return s.request(ctx, "addPost", in)
}
