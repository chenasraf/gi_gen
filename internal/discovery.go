package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chenasraf/gi_gen/internal/utils"
	"golang.org/x/exp/maps"
)

func AutoDiscover(allFiles []string) ([]string, map[string]string) {
	list := discoverByExplicitProjectType()

	if len(list) == 0 {
		list = discoverByExistingPatterns(allFiles)
	}

	return maps.Keys(list), list
}

func readFromSelections(allFiles []string) ([]string, map[string]string) {
	answer := askDiscovery()

	if !answer {
		return readAllFiles(allFiles)
	}

	return AutoDiscover(allFiles)
}

func readAllFiles(allFiles []string) ([]string, map[string]string) {
	baseNames := []string{}
	for _, fn := range allFiles {
		basename := filepath.Base(fn)
		langName := basename[:strings.Index(basename, ".")]
		baseNames = append(baseNames, langName)
	}
	return baseNames, getAllFiles(allFiles)
}

func getAllFiles(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := utils.ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		files[langName] = contents
	}

	return files
}

func discoverByExistingPatterns(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := utils.ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		if findPatternFileMatches(contents) {
			files[langName] = contents
		}
	}
	return files
}

func discoverByExplicitProjectType() map[string]string {
	wd, err := os.Getwd()
	utils.HandleErr(err)

	discoveryMap := make(map[string]string)

	// Common workspace files
	discoveryMap["app/manifests/AndroidManifest.xml"] = "Android"
	discoveryMap["composer.json"] = "Composer"
	discoveryMap["pubspec.ya?ml"] = "Dart"
	discoveryMap["go.{mod,sum}"] = "Go"
	discoveryMap["_config.ya?ml"] = "Jekyll"
	discoveryMap["Jenkinsfile"] = "JENKINS_HOME"
	discoveryMap["jobs"] = "JENKINS_HOME"
	discoveryMap["package.json"] = "Node"
	discoveryMap["[tj]sconfig.json"] = "Node"
	discoveryMap["Gemfile"] = "Ruby"
	discoveryMap["bin/rails"] = "Rails"
	discoveryMap["config/boot.rb"] = "Rails"
	discoveryMap["*.uproject"] = "UnrealEngine"

	// Extensions
	discoveryMap["*.as"] = "Actionscript"
	discoveryMap["*.{adb,ada,ads}"] = "Ada"
	discoveryMap["*.agda"] = "Agda"
	discoveryMap["*.{c,cats,h,idc,w}"] = "C"
	discoveryMap["*.{cpp,c++,cc,cp,cxx,h,h++,hh,hpp,hxx,inc,inl,ipp,tcc,tpp}"] = "C++"
	discoveryMap["*.{cmake,cmake.in}"] = "CMake"
	discoveryMap["*.{clj,boot,cl2,cljc,cljs,cljs.hl,cljscm,cljx,hic}"] = "Clojure"
	discoveryMap["*.{coq,v}"] = "Coq"
	discoveryMap["*.{cu,cuh}"] = "CUDA"
	discoveryMap["*.{d,di}"] = "D"
	discoveryMap["*.dm"] = "DM"
	discoveryMap["*.dart"] = "Dart"
	discoveryMap["*.{sch,brd}"] = "Eagle"
	discoveryMap["*.{ex,exs}"] = "Elixir"
	discoveryMap["*.elm"] = "Elm"
	discoveryMap["*.{erl,es,escript,hrl,xrl,yrl}"] = "Erlang"
	discoveryMap["*.{f90,f,f03,f08,f77,f95,for,fpp}"] = "Fortran"
	discoveryMap["*.{fy,fancypack}"] = "Fancy"
	discoveryMap["*.go"] = "Go"
	discoveryMap["*.godot"] = "Godot"
	discoveryMap["*.gradle"] = "Gradle"
	discoveryMap["*.{hs,hsc}"] = "Haskell"
	discoveryMap["*.{idr,lidr}"] = "Idris"
	discoveryMap["*.java"] = "Java"
	discoveryMap["*.jl"] = "Julia"
	discoveryMap["*.{sch,brd,kicad_pcb}"] = "KiCad"
	discoveryMap["*.{kt,ktm,kts}"] = "Kotlin"
	discoveryMap["*.lvproj"] = "LabVIEW"
	discoveryMap["*.{ly,ily}"] = "Lilypond"
	discoveryMap["*.{lua,fcgi,nse,pd_lua,rbxs,wlua}"] = "Lua"
	discoveryMap["*.{m,moo}"] = "Mercury"
	discoveryMap["*.js"] = "Node"
	discoveryMap["*.{ml,eliom,eliomi,ml4,mli,mll,mly}"] = "OCaml"
	discoveryMap["*.{m,h}"] = "Objective-C"
	discoveryMap["*.opa"] = "Opa"
	discoveryMap["*.{pl,al,cgi,fcgi,perl,ph,plx,pm,pod,psgi,t}"] = "Perl"
	discoveryMap["*.pde"] = "Processing"
	discoveryMap["*.purs"] = "PureScript"
	discoveryMap["*.{py,bzl,cgi,fcgi,gyp,lmi,pyde,pyp,pyt,pyw,rpy,tac,wsgi,xpy}"] = "Python"
	discoveryMap["*.{r,rd,rsx}"] = "R"
	discoveryMap["*.{rkt,rktd,rktl,scrbl}"] = "Racket"
	discoveryMap["*.{rb,builder,fcgi,gemspec,god,irbrc,jbuilder,mspec,pluginspec,podspec,rabl,rake,rbuild,rbw,rbx,ru,ruby,thor,watchr}"] = "Ruby"
	discoveryMap["*.{rs,rs.in}"] = "Rust"
	discoveryMap["*.sass"] = "Sass"
	discoveryMap["*.{scala,sbt,sc}"] = "Scala"
	discoveryMap["*.{scm,sld,sls,sps,ss}"] = "Scheme"
	discoveryMap["*.{st,cs}"] = "Smalltalk"
	discoveryMap["*.swift"] = "Swift"
	discoveryMap["*.{tex,aux,bbx,bib,cbx,cls,dtx,ins,lbx,ltx,mkii,mkiv,mkvi,sty,toc}"] = "TeX"
	discoveryMap["*.unity"] = "Unity"
	discoveryMap["*.{xojo_code,xojo_menu,xojo_report,xojo_script,xojo_toolbar,xojo_window}"] = "Xojo"
	discoveryMap["*.zep"] = "Zephir"

	results := make(map[string]string)

	for _, key := range maps.Keys(discoveryMap) {
		langName := discoveryMap[key]
		ignoreFile := filepath.Join(GetCacheDir(), langName+".gitignore")
		checkFile := filepath.Join(wd, key)

		_, keyExists := results[langName]
		if !keyExists && utils.GlobExists(checkFile) {
			results[langName] = utils.ReadFile(ignoreFile)
		}
	}

	return results
}
