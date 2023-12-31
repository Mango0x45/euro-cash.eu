import logging
from http import HTTPMethod, HTTPStatus
from typing import NamedTuple

import flask
import flask_babel
from flask import Blueprint, Response
from flask_babel import _
from selectolax.parser import HTMLParser

import util
import watchers.articles
from xtypes import LOCALES, Cookie, Theme


class Article(NamedTuple):
	headline: str
	iso_8601: str
	date: str
	sample: str


root = Blueprint("root", __name__, url_prefix="/")


@root.route("/", methods=[HTTPMethod.GET])
def index() -> str:
	try:
		newest = watchers.articles.watcher.articles[0]
	except IndexError:
		return flask.render_template("index.html", article=None)

	try:
		path = util.from_root(f"templates/news/articles/{newest.date}.html")
		with open(path, "r") as f:
			p = HTMLParser(f.read())
	except Exception as e:
		logging.root.error(e)
		p = None

	if p is not None and (sample := p.css_first("main p")):
		sample = util.strip_jinja(sample.text())
	else:
		sample = _("No article preview found")

	newest = Article(
		headline=newest.headline,
		iso_8601=newest.date.isoformat(),
		date=flask_babel.format_date(newest.date, "short"),
		sample=sample,
	)
	return flask.render_template("index.html", article=newest)


@root.route("/jargon", methods=[HTTPMethod.GET])
def jargon() -> str:
	return flask.render_template("jargon.html")


@root.route("/about", methods=[HTTPMethod.GET])
def about() -> str:
	return flask.render_template("about.html")


@root.route("/set-language", methods=[HTTPMethod.GET, HTTPMethod.POST])
def set_language() -> Response | tuple[str, HTTPStatus]:
	r = flask.request

	try:
		loc = r.form[Cookie.LOCALE]
		assert loc in LOCALES
	except KeyError:
		s = f"‘{Cookie.LOCALE}’ key missing from request form"
		logging.root.warning(s)
		return s, HTTPStatus.BAD_REQUEST
	except AssertionError:
		s = f"Locale ‘{loc}’ is not an available locale"  # type: ignore
		logging.root.warning(s)
		return s, HTTPStatus.BAD_REQUEST

	url = r.cookies.get(Cookie.REDIRECT, default="/")
	resp = flask.make_response(flask.redirect(url))
	resp.delete_cookie(Cookie.REDIRECT)
	resp.set_cookie(key=Cookie.LOCALE, value=loc, max_age=2**31 - 1)
	logging.root.debug(f"‘{Cookie.LOCALE}’ cookie set to ‘{loc}’")

	return resp


@root.route("/toggle-theme", methods=[HTTPMethod.POST])
def toggle_theme() -> Response:
	r = flask.request

	theme = r.cookies.get(Cookie.THEME)
	resp = flask.make_response(flask.redirect(r.referrer))
	theme = Theme.DARK if theme == Theme.LIGHT else Theme.LIGHT
	resp.set_cookie(
		key=Cookie.THEME,
		value=theme,
		max_age=2**31 - 1,
	)
	logging.root.debug(f"‘{Cookie.THEME}’ cookie set to ‘{theme}’")

	return resp
