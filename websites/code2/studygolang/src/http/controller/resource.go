// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"html/template"
	"logic"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"

	"model"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeResource, logic.ResourceComment{})
	// service.RegisterLikeObject(model.TYPE_RESOURCE, service.ResourceLike{})
}

type ResourceController struct{}

// 注册路由
func (this *ResourceController) RegisterRoute(e *echo.Echo) {
	e.Get("/resources", echo.HandlerFunc(this.ReadList))
	e.Get("/resources/cat/:catid", echo.HandlerFunc(this.ReadCatResources))
	e.Get("/resources/:id", echo.HandlerFunc(this.Detail))
	e.Any("/resources/new", echo.HandlerFunc(this.Create))
	e.Any("/resources/modify", echo.HandlerFunc(this.Modify))
}

// ReadList 资源索引页
func (ResourceController) ReadList(ctx echo.Context) error {
	return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
}

// ReadCatResources 某个分类的资源列表
func (ResourceController) ReadCatResources(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.Query("p"), 1)
	paginator := logic.NewPaginator(curPage)
	catid := goutils.MustInt(ctx.Param("catid"))

	resources, total := logic.DefaultResource.FindByCatid(ctx, paginator, catid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(Request(ctx).URL.Path)

	return render(ctx, "resources/index.html", map[string]interface{}{"activeResources": "active", "resources": resources, "categories": logic.AllCategory, "page": template.HTML(pageHtml), "curCatid": catid})
}

// Detail 某个资源详细页
func (ResourceController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id > 0 {
		return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
	}
	resource, comments := logic.DefaultResource.FindById(ctx, id)
	if len(resource) == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
	}

	likeFlag := 0
	hadCollect := 0
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		id := resource["id"].(int)
		likeFlag = logic.DefaultLike.HadLike(ctx, me.Uid, id, model.TypeResource)
		hadCollect = logic.DefaultFavorite.HadFavorite(ctx, me.Uid, id, model.TypeResource)
	}

	// service.Views.Incr(req, model.TYPE_RESOURCE, util.MustInt(vars["id"]))

	return render(ctx, "resources/detail.html,common/comment.html", map[string]interface{}{"activeResources": "active", "resource": resource, "comments": comments, "likeflag": likeFlag, "hadcollect": hadCollect})
}

// Create 发布新资源
func (ResourceController) Create(ctx echo.Context) error {
	title := ctx.Form("title")
	// 请求新建资源页面
	if title == "" || Request(ctx).Method != "POST" {
		return render(ctx, "resources/new.html", map[string]interface{}{"activeResources": "active", "categories": logic.AllCategory})
	}

	errMsg := ""
	resForm := ctx.Form("form")
	if resForm == model.LinkForm {
		if ctx.Form("url") == "" {
			errMsg = "url不能为空"
		}
	} else {
		if ctx.Form("content") == "" {
			errMsg = "内容不能为空"
		}
	}
	if errMsg != "" {
		return fail(ctx, 1, errMsg)
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultResource.Publish(ctx, me, Request(ctx).Form)
	if err != nil {
		return fail(ctx, 2, "内部服务错误，请稍候再试！")
	}

	return success(ctx, nil)
}

// Modify 修改資源
func (ResourceController) Modify(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Form("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
	}

	// 请求编辑資源页面
	if Request(ctx).Method != "POST" {
		resource := logic.DefaultResource.FindResource(ctx, id)
		return render(ctx, "resources/new.html", map[string]interface{}{"resource": resource, "activeResources": "active", "categories": service.AllCategory})
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultResource.Publish(ctx, me, Request(ctx).Form)
	if err != nil {
		if err == service.NotModifyAuthorityErr {
			return ctx.Error("没有权限修改")
		}
		return fail(ctx, 2, "内部服务错误，请稍候再试！")
	}

	return success(ctx, nil)
}