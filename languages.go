package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chenasraf/utils"
)

// Attempts to auto-select languages.
//
// May return more than 1 matching language, or none. In both cases the choice should be delegated
// to the user, either with the returned candidates, or with the entire list if none were found.
//
// Returns language names, and corresponding matches
func getAutoSelectCandidates() ([]string, map[string]LanguageMatch) {
	languages, matches := getLanguagesByGlobs()

	if len(matches) > 0 {
		return languages, matches
	}

	if len(matches) == 0 {
		allTemplates := getAllIgnoreTemplates()
		languages, matches = discoverByExistingPatterns(allTemplates)
	}

	return languages, matches
}

// Try to find files in the current directory that match any glob patterns for languages by
// iterating over each language's gitignore, and attempting to find a match in there.
//
// Unlike `getLanguagesByGlobs`, this does not use a static list, but uses each of the language
// template files as base.
func discoverByExistingPatterns(ignoreFiles map[string]string) ([]string, map[string]LanguageMatch) {
	matches := make(map[string]LanguageMatch)

	for ignoreFile, contents := range ignoreFiles {
		filename := filepath.Base(ignoreFile)
		langName := filename[:strings.Index(filename, ".")]

		if found, pattern := findPatternFileMatches(contents); found {
			matches[langName] = LanguageMatch{
				language:       langName,
				matchedPattern: pattern,
				contents:       contents,
			}
		}
	}
	return utils.MapKeys(matches), matches
}

// Tries to find files in the current directory that match any glob patterns for languages (static discovery).
//
// e.g. matches package.json to Node.js, pubspec.yaml to Dart, etc
//
// @returns map of language name to LanguageMatch
func getLanguagesByGlobs() ([]string, map[string]LanguageMatch) {
	patternMap := getLanguagePatterns()
	patternsList := utils.MapKeys(patternMap)

	matches := make(map[string]LanguageMatch)

	wd, err := os.Getwd()
	utils.HandleErr(err)

	for _, globPart := range patternsList {
		langName := patternMap[globPart]
		ignoreFile := filepath.Join(GetCacheDir(), langName+".gitignore")
		projectDirGlob := filepath.Join(wd, globPart)

		_, isAlreadyMatched := matches[langName]
		if !isAlreadyMatched && utils.GlobExists(projectDirGlob) {
			contents := utils.ReadFile(ignoreFile)
			matches[langName] = LanguageMatch{
				language:       langName,
				matchedPattern: globPart,
				contents:       contents,
			}
		}
	}

	keys := utils.MapKeys(matches)
	return keys, matches
}

// returns all gitignore templates with their contents
func getAllIgnoreTemplates() map[string]string {
	allTemplates, err := InitCache()
	utils.HandleErr(err)

	files := make(map[string]string)

	for _, filename := range allTemplates {
		contents := utils.ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		files[langName] = contents
	}

	return files
}

// Returns a map of glob pattrns to respective language names
func getLanguagePatterns() map[string]string {
	discoveryMap := make(map[string]string)

	// Common workspace files
	discoveryMap["*.uproject"] = "UnrealEngine"
	discoveryMap["Gemfile"] = "Ruby"
	discoveryMap["Jenkinsfile"] = "JENKINS_HOME"
	discoveryMap["[tj]sconfig.json"] = "Node"
	discoveryMap["_config.ya?ml"] = "Jekyll"
	discoveryMap["app/manifests/AndroidManifest.xml"] = "Android"
	discoveryMap["bin/rails"] = "Rails"
	discoveryMap["composer.json"] = "Composer"
	discoveryMap["config/boot.rb"] = "Rails"
	discoveryMap["go.{mod,sum}"] = "Go"
	discoveryMap["jobs"] = "JENKINS_HOME"
	discoveryMap["package.json"] = "Node"
	discoveryMap["pubspec.ya?ml"] = "Dart"
	discoveryMap["configure.ac"] = "Autotools"

	// Extensions
	discoveryMap["**/*.agda"] = "Agda"
	discoveryMap["**/*.al"] = "AL"
	discoveryMap["**/*.as"] = "Actionscript"
	discoveryMap["**/*.dart"] = "Dart"
	discoveryMap["**/*.dm"] = "DM"
	discoveryMap["**/*.el"] = "Elisp"
	discoveryMap["**/*.elm"] = "Elm"
	discoveryMap["**/*.gcov"] = "Gcov"
	discoveryMap["**/*.go"] = "Go"
	discoveryMap["**/*.godot"] = "Godot"
	discoveryMap["**/*.gradle"] = "Gradle"
	discoveryMap["**/*.java"] = "Java"
	discoveryMap["**/*.jl"] = "Julia"
	discoveryMap["**/*.js"] = "Node"
	discoveryMap["**/*.lvproj"] = "LabVIEW"
	discoveryMap["**/*.nim"] = "Nim"
	discoveryMap["**/*.opa"] = "Opa"
	discoveryMap["**/*.pde"] = "Processing"
	discoveryMap["**/*.rb"] = "ChefCookbook"
	discoveryMap["**/*.sass"] = "Sass"
	discoveryMap["**/*.swift"] = "Swift"
	discoveryMap["**/*.unity"] = "Unity"
	discoveryMap["**/*.zep"] = "Zephir"
	discoveryMap["**/*.{adb,ada,ads}"] = "Ada"
	discoveryMap["**/*.{c,cats,h,idc,w}"] = "C"
	discoveryMap["**/*.{c,h}"] = "Sdcc"
	discoveryMap["**/*.{cfm,cfc}"] = "CFWheels"
	discoveryMap["**/*.{clj,boot,cl2,cljc,cljs,cljs.hl,cljscm,cljx,hic}"] = "Clojure"
	discoveryMap["**/*.{clj,cljs,edn}"] = "Leiningen"
	discoveryMap["**/*.{cls,trigger,page,component}"] = "ForceDotCom"
	discoveryMap["**/*.{cmake,cmake.in}"] = "CMake"
	discoveryMap["**/*.{coq,v}"] = "Coq"
	discoveryMap["**/*.{cpp,c++,cc,cp,cxx,h,h++,hh,hpp,hxx,inc,inl,ipp,tcc,tpp}"] = "C++"
	discoveryMap["**/*.{cpp,h}"] = "FlaxEngine"
	discoveryMap["**/*.{cpp,h}"] = "Qt"
	discoveryMap["**/*.{cpp,h}"] = "ROS"
	discoveryMap["**/*.{cpp,h}"] = "TwinCAT3"
	discoveryMap["**/*.{ctp,php}"] = "CakePHP"
	discoveryMap["**/*.{cu,cuh}"] = "CUDA"
	discoveryMap["**/*.{d,di}"] = "D"
	discoveryMap["**/*.{erl,es,escript,hrl,xrl,yrl}"] = "Erlang"
	discoveryMap["**/*.{ex,exs}"] = "Elixir"
	discoveryMap["**/*.{f90,f,f03,f08,f77,f95,for,fpp}"] = "Fortran"
	discoveryMap["**/*.{fmb,pll,mmb}"] = "OracleForms"
	discoveryMap["**/*.{fy,fancypack}"] = "Fancy"
	discoveryMap["**/*.{groovy,java}"] = "Grails"
	discoveryMap["**/*.{hs,hsc}"] = "Haskell"
	discoveryMap["**/*.{idr,lidr}"] = "Idris"
	discoveryMap["**/*.{java,kt,xml}"] = "Android"
	discoveryMap["**/*.{java,scala}"] = "PlayFramework"
	discoveryMap["**/*.{java,xml}"] = "GWT"
	discoveryMap["**/*.{java,xml}"] = "SeamGen"
	discoveryMap["**/*.{js,ts,xml}"] = "AppceleratorTitanium"
	discoveryMap["**/*.{js,ts}"] = "ExtJs"
	discoveryMap["**/*.{json,hcl}"] = "Packer"
	discoveryMap["**/*.{kt,ktm,kts}"] = "Kotlin"
	discoveryMap["**/*.{lisp,lsp}"] = "CommonLisp"
	discoveryMap["**/*.{lua,fcgi,nse,pd_lua,rbxs,wlua}"] = "Lua"
	discoveryMap["**/*.{ly,ily}"] = "Lilypond"
	discoveryMap["**/*.{m,h}"] = "Objective-C"
	discoveryMap["**/*.{m,moo}"] = "Mercury"
	discoveryMap["**/*.{md,html}"] = "GitBook"
	discoveryMap["**/*.{ml,eliom,eliomi,ml4,mli,mll,mly}"] = "OCaml"
	discoveryMap["**/*.{mus,xml}"] = "Finale"
	discoveryMap["**/*.{pas,dpr,dpk,dfm,dproj}"] = "Delphi"
	discoveryMap["**/*.{php,html,css,js}"] = "CodeIgniter"
	discoveryMap["**/*.{php,html,css,js}"] = "Concrete5"
	discoveryMap["**/*.{php,html,css,js}"] = "CraftCMS"
	discoveryMap["**/*.{php,html,css,js}"] = "Drupal"
	discoveryMap["**/*.{php,html,css,js}"] = "ExpressionEngine"
	discoveryMap["**/*.{php,html,css,js}"] = "FuelPHP"
	discoveryMap["**/*.{php,html,css,js}"] = "Joomla"
	discoveryMap["**/*.{php,html,css,js}"] = "Kohana"
	discoveryMap["**/*.{php,html,css,js}"] = "Laravel"
	discoveryMap["**/*.{php,html,css,js}"] = "LemonStand"
	discoveryMap["**/*.{php,html,css,js}"] = "Lithium"
	discoveryMap["**/*.{php,html,css,js}"] = "Magento"
	discoveryMap["**/*.{php,html,css,js}"] = "MetaProgrammingSystem"
	discoveryMap["**/*.{php,html,css,js}"] = "Nanoc"
	discoveryMap["**/*.{php,html,css,js}"] = "OpenCart"
	discoveryMap["**/*.{php,html,css,js}"] = "Phalcon"
	discoveryMap["**/*.{php,html,css,js}"] = "Plone"
	discoveryMap["**/*.{php,html,css,js}"] = "Prestashop"
	discoveryMap["**/*.{php,html,css,js}"] = "Qooxdoo"
	discoveryMap["**/*.{php,html,css,js}"] = "Raku"
	discoveryMap["**/*.{php,html,css,js}"] = "Scrivener"
	discoveryMap["**/*.{php,html,css,js}"] = "Stella"
	discoveryMap["**/*.{php,html,css,js}"] = "SugarCRM"
	discoveryMap["**/*.{php,html,css,js}"] = "Symfony"
	discoveryMap["**/*.{php,html,css,js}"] = "SymphonyCMS"
	discoveryMap["**/*.{php,html,css,js}"] = "Textpattern"
	discoveryMap["**/*.{php,html,css,js}"] = "Typo3"
	discoveryMap["**/*.{php,html,css,js}"] = "VVVV"
	discoveryMap["**/*.{php,html,css,js}"] = "WordPress"
	discoveryMap["**/*.{php,html,css,js}"] = "Yeoman"
	discoveryMap["**/*.{php,html,css,js}"] = "Yii"
	discoveryMap["**/*.{php,html,css,js}"] = "ZendFramework"
	discoveryMap["**/*.{pkgbuild,install}"] = "ArchLinuxPackages"
	discoveryMap["**/*.{pl,al,cgi,fcgi,perl,ph,plx,pm,pod,psgi,t}"] = "Perl"
	discoveryMap["**/*.{pxp,ipf}"] = "IGORPro"
	discoveryMap["**/*.{py,bzl,cgi,fcgi,gyp,lmi,pyde,pyp,pyt,pyw,rpy,tac,wsgi,xpy}"] = "Python"
	discoveryMap["**/*.{py,html,css,js}"] = "TurboGears2"
	discoveryMap["**/*.{r,rd,rsx}"] = "R"
	discoveryMap["**/*.{rb,builder,fcgi,gemspec,god,irbrc,jbuilder,mspec,pluginspec,podspec,rabl,rake,rbuild,rbw,rbx,ru,ruby,thor,watchr}"] = "Ruby"
	discoveryMap["**/*.{rb,erb}"] = "RhodesRhomobile"
	discoveryMap["**/*.{rkt,rktd,rktl,scrbl}"] = "Racket"
	discoveryMap["**/*.{rs,rs.in}"] = "Rust"
	discoveryMap["**/*.{scala,sbt,sc}"] = "Scala"
	discoveryMap["**/*.{sch,brd,kicad_pcb}"] = "KiCad"
	discoveryMap["**/*.{sch,brd}"] = "Eagle"
	discoveryMap["**/*.{scm,sld,sls,sps,ss}"] = "Scheme"
	discoveryMap["**/*.{skp,rb}"] = "SketchUp"
	discoveryMap["**/*.{sln,csproj}"] = "VisualStudio"
	discoveryMap["**/*.{st,cs}"] = "Smalltalk"
	discoveryMap["**/*.{tex,aux,bbx,bib,cbx,cls,dtx,ins,lbx,ltx,mkii,mkiv,mkvi,sty,toc}"] = "TeX"
	discoveryMap["**/*.{tf,hcl}"] = "Terraform"
	discoveryMap["**/*.{wscript,py}"] = "Waf"
	discoveryMap["**/*.{xml,java}"] = "Maven"
	discoveryMap["**/*.{xojo_code,xojo_menu,xojo_report,xojo_script,xojo_toolbar,xojo_window}"] = "Xojo"
	discoveryMap["**/*.{yaml,yml}"] = "AppEngine"

	return discoveryMap
}
