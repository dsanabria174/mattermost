// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model


func (o *PostPatch) IsOnlyRemoveLinkPreview() bool {
	if o.Props == nil || o.Message != nil || o.FileIds != nil || o.IsPinned != nil || o.HasReactions != nil {
		return false
	}
	return (*o.Props)[PostPropsRemoveLinkPreview] == "true"
}
