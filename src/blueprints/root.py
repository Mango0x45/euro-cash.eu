import os
import re
from http import HTTPMethod, HTTPStatus
from typing import NamedTuple

import flask
import flask_babel
from flask import Blueprint, Response

import article_watcher
from config import Cookie, Locale, Theme
from util import _


class Article(NamedTuple):
	headline: str
	iso_8601: str
	date: str
	sample: str


root = Blueprint("root", __name__, url_prefix="/")


@root.route("/", methods=[HTTPMethod.GET])
def index() -> str:
	try:
		newest = article_watcher.watcher.articles[0]
	except IndexError:
		return flask.render_template("index.html", article=None)

	path = os.path.dirname(__file__)
	path = os.path.join(path, f"../templates/news/articles/{newest.date}.html")
	with open(path, "r") as f:
		txt = f.read()

	ms = re.findall(r'<p>[^"]*"([^"]*)"[^<]*</p>', txt, flags=re.DOTALL)
	try:
		print(ms)
		sample: str = ms[1]
	except IndexError:
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
		assert loc in Locale
	except KeyError:
		return (
			f'"{Cookie.LOCALE}" key missing from request form',
			HTTPStatus.BAD_REQUEST,
		)
	except AssertionError:
		return f'Locale "{loc}" is not an available locale', HTTPStatus.BAD_REQUEST  # type: ignore

	url = r.cookies.get(Cookie.REDIRECT, default="/")
	resp = flask.make_response(flask.redirect(url))
	resp.delete_cookie(Cookie.REDIRECT)
	resp.set_cookie(key=Cookie.LOCALE, value=loc, max_age=2**31 - 1)

	return resp


@root.route("/toggle-theme", methods=[HTTPMethod.POST])
def toggle_theme() -> Response:
	r = flask.request

	theme = r.cookies.get(Cookie.THEME)
	resp = flask.make_response(flask.redirect(r.referrer))
	resp.set_cookie(
		key=Cookie.THEME,
		value=Theme.DARK if theme == Theme.LIGHT else Theme.LIGHT,
		max_age=2**31 - 1,
	)

	return resp
