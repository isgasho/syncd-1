// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
)

func SpaceNew(c *goweb.Context) error {
    return spaceUpdate(c, 0)
}

func SpaceEdit(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
        return syncd.RenderParamError("id can not be empty")
    }
    return spaceUpdate(c, id)
}

func spaceUpdate(c *goweb.Context, id int) error {
    name := c.PostForm("name")
    if name == "" {
        return syncd.RenderParamError("name can not be empty")
    }
    var space *projectService.Space
    space = &projectService.Space{
        ID: id,
        Name: name,
    }
    exists, err := space.CheckExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderAppError("space data update failed, space name have exists")
    }
    space = &projectService.Space{
        ID: id,
        Name: name,
        Description: c.PostForm("description"),
    }
    if err := space.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func SpaceList(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.GetInt("limit"), c.Query("keyword")
    space := &projectService.Space{}
    list, total, err := space.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    //check if project exists in the space
    var newList []map[string]interface{}
    for _, l := range list {
        project := &projectService.Project{
            SpaceId: l.ID,
        }
        exists, err := project.CheckSpaceHaveProject()
        if err != nil {
            return syncd.RenderAppError(err.Error())
        }
        newList = append(newList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "description": l.Description,
            "have_project": exists,
            "ctime": l.Ctime,
        })
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": newList,
        "total": total,
    })
}

func SpaceDetail(c *goweb.Context) error {
    space := &projectService.Space{
        ID: c.QueryInt("id"),
    }
    if err := space.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, space)
}

func SpaceDelete(c *goweb.Context) error {
    var (
        id int
        exists bool
        err error
    )
    id = c.PostFormInt("id")
    proj := &projectService.Project{
        SpaceId: id,
    }
    exists, err = proj.CheckSpaceHaveProject()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderAppError("space delete failed, project in space is not empty")
    }
    space := &projectService.Space{
        ID: id,
    }
    if err = space.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func SpaceExists(c *goweb.Context) error {
    keyword, id := c.Query("keyword"), c.QueryInt("id")
    if keyword == "" {
        return syncd.RenderParamError("params error")
    }
    space := &projectService.Space{
        ID: id,
        Name: keyword,
    }
    exists, err := space.CheckExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "exists": exists,
    })
}

