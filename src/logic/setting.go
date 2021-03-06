// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"encoding/json"
	"net/url"
	"strings"

	"model"

	"github.com/polaris1119/goutils"
	"golang.org/x/net/context"
)

type SettingLogic struct{}

var DefaultSetting = SettingLogic{}

func (SettingLogic) Update(ctx context.Context, form url.Values) error {
	objLog := GetLogger(ctx)

	name := form.Get("name")
	if name != "" {
		WebsiteSetting.Name = name
	}

	domain := form.Get("domain")
	if domain != "" {
		WebsiteSetting.Domain = domain
	}

	titleSuffix := form.Get("title_suffix")
	if titleSuffix != "" {
		WebsiteSetting.TitleSuffix = titleSuffix
	}

	favicon := form.Get("favicon")
	if favicon != "" {
		WebsiteSetting.Favicon = favicon
	}

	startYear := goutils.MustInt(form.Get("start_year"))
	if startYear != 0 {
		WebsiteSetting.StartYear = startYear
	}

	logo := form.Get("logo")
	if logo != "" {
		WebsiteSetting.Logo = logo
	}

	WebsiteSetting.BlogUrl = form.Get("blog_url")

	slogan := form.Get("slogan")
	if slogan != "" {
		WebsiteSetting.Slogan = slogan
	}

	WebsiteSetting.Beian = form.Get("beian")

	WebsiteSetting.ReadingMenu = form.Get("reading_menu")

	if docNameSlice, ok := form["doc_name"]; ok {
		docUrlSlice := form["doc_url"]

		docMenus := make([]*model.DocMenu, len(docNameSlice))
		for i, docName := range docNameSlice {
			docMenus[i] = &model.DocMenu{
				Name: docName,
				Url:  docUrlSlice[i],
			}
		}

		docMenusBytes, err := json.Marshal(docMenus)
		if err != nil {
			objLog.Errorln("marshal doc menu error:", err)
			return err
		}

		WebsiteSetting.DocMenus = docMenus
		WebsiteSetting.DocsMenu = string(docMenusBytes)
	}

	if navNameSlice, ok := form["nav_name"]; ok {
		navUrlSlice := form["nav_url"]

		footerNavs := make([]*model.FooterNav, len(navNameSlice))
		for i, navName := range navNameSlice {
			outerWeb := true
			if strings.HasPrefix(navUrlSlice[i], "/") {
				outerWeb = false
			}
			footerNavs[i] = &model.FooterNav{
				Name:      navName,
				Url:       navUrlSlice[i],
				OuterSite: outerWeb,
			}
		}

		footerNavsBytes, err := json.Marshal(footerNavs)
		if err != nil {
			objLog.Errorln("marshal footer nav error:", err)
			return err
		}

		WebsiteSetting.FooterNavs = footerNavs
		WebsiteSetting.FooterNav = string(footerNavsBytes)
	}

	if frLogoImageSlice, ok := form["fr_logo_image"]; ok {
		frLogoUrlSlice := form["fr_logo_url"]
		frWidthSlice := form["fr_logo_width"]
		frHeightSlice := form["fr_logo_height"]

		friendLogos := make([]*model.FriendLogo, len(frLogoImageSlice))
		for i, frLogoImage := range frLogoImageSlice {
			friendLogos[i] = &model.FriendLogo{
				Image:  frLogoImage,
				Url:    frLogoUrlSlice[i],
				Width:  frWidthSlice[i],
				Height: frHeightSlice[i],
			}
		}

		friendLogosBytes, err := json.Marshal(friendLogos)
		if err != nil {
			objLog.Errorln("marshal friend logo error:", err)
			return err
		}

		WebsiteSetting.FriendLogos = friendLogos
		WebsiteSetting.FriendsLogo = string(friendLogosBytes)
	}

	_, err := MasterDB.Update(WebsiteSetting)
	if err != nil {
		objLog.Errorln("Update setting error:", err)
		return err
	}

	return nil
}
