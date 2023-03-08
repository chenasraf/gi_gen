use std::{collections::HashMap, fs, io::Error, path::PathBuf};

pub fn get_language_candidates(path: &PathBuf) -> Result<Vec<String>, Error> {
    let mut result: Vec<String> = Vec::new();
    let patterns = get_glob_patterns();

    for (pattern, lang) in &patterns {
        let files: Vec<Result<fs::DirEntry, Error>> = path.read_dir()?.collect();
        let glob_compiler = match globset::Glob::new(pattern) {
            Ok(it) => it,
            Err(err) => return Err(Error::new(std::io::ErrorKind::Other, err)),
        };
        let glob_pattern = glob_compiler.compile_matcher();
        for file in files {
            let file = file?;
            let file_name = match file.file_name().into_string() {
                Ok(name) => name,
                Err(_) => Err(Error::new(
                    std::io::ErrorKind::Other,
                    "Could not convert file name to string",
                ))?,
            };
            if glob_pattern.is_match(&file_name) {
                result.push(lang.to_string());
            }
        }
    }
    result.sort();

    Ok(result)
}

fn get_glob_patterns() -> HashMap<String, String> {
    let mut result: HashMap<String, String> = HashMap::new();
    // Common workspace files
    result.insert(
        String::from("app/manifests/AndroidManifest.xml"),
        String::from("Android"),
    );
    result.insert(String::from("composer.json"), String::from("Composer"));
    result.insert(String::from("pubspec.{yml,yaml}"), String::from("Dart"));
    result.insert(String::from("go.{mod,sum}"), String::from("Go"));
    result.insert(String::from("_config.{yml,yaml}"), String::from("Jekyll"));
    result.insert(String::from("Jenkinsfile"), String::from("JENKINS_HOME"));
    result.insert(String::from("jobs"), String::from("JENKINS_HOME"));
    result.insert(String::from("package.json"), String::from("Node"));
    result.insert(String::from("[tj]sconfig.json"), String::from("Node"));
    result.insert(String::from("Gemfile"), String::from("Ruby"));
    result.insert(String::from("bin/rails"), String::from("Rails"));
    result.insert(String::from("config/boot.rb"), String::from("Rails"));
    result.insert(String::from("*.uproject"), String::from("UnrealEngine"));
    result.insert(String::from("Cargo.toml"), String::from("Rust"));

    // Extensions
    result.insert(String::from("*.as"), String::from("Actionscript"));
    result.insert(String::from("*.{adb,ada,ads}"), String::from("Ada"));
    result.insert(String::from("*.agda"), String::from("Agda"));
    result.insert(String::from("*.{c,cats,h,idc,w}"), String::from("C"));
    result.insert(
        String::from("*.{cpp,c++,cc,cp,cxx,h,h++,hh,hpp,hxx,inc,inl,ipp,tcc,tpp}"),
        String::from("C++"),
    );
    result.insert(String::from("*.{cmake,cmake.in}"), String::from("CMake"));
    result.insert(
        String::from("*.{clj,boot,cl2,cljc,cljs,cljs.hl,cljscm,cljx,hic}"),
        String::from("Clojure"),
    );
    result.insert(String::from("*.{coq,v}"), String::from("Coq"));
    result.insert(String::from("*.{cu,cuh}"), String::from("CUDA"));
    result.insert(String::from("*.{d,di}"), String::from("D"));
    result.insert(String::from("*.dm"), String::from("DM"));
    result.insert(String::from("*.dart"), String::from("Dart"));
    result.insert(String::from("*.{sch,brd}"), String::from("Eagle"));
    result.insert(String::from("BUILD.gn"), String::from("Electron"));
    result.insert(String::from("*.{ex,exs}"), String::from("Elixir"));
    result.insert(String::from("*.elmi?"), String::from("Elm"));
    result.insert(
        String::from("*.{erl,es,escript,hrl,xrl,yrl}"),
        String::from("Erlang"),
    );
    result.insert(
        String::from("*.{f90,f,f03,f08,f77,f95,for,fpp}"),
        String::from("Fortran"),
    );
    result.insert(String::from("*.{fy,fancypack}"), String::from("Fancy"));
    result.insert(String::from("*.go"), String::from("Go"));
    result.insert(String::from("*.godot"), String::from("Godot"));
    result.insert(String::from("*.gradle"), String::from("Gradle"));
    result.insert(String::from("*.{hs,hsc}"), String::from("Haskell"));
    result.insert(String::from("*.{idr,lidr}"), String::from("Idris"));
    result.insert(String::from("*.java"), String::from("Java"));
    result.insert(String::from("*.jl"), String::from("Julia"));
    result.insert(String::from("*.{sch,brd,kicad_pcb}"), String::from("KiCad"));
    result.insert(String::from("*.{kt,ktm,kts}"), String::from("Kotlin"));
    result.insert(String::from("*.lvproj"), String::from("LabVIEW"));
    result.insert(String::from("*.{ly,ily}"), String::from("Lilypond"));
    result.insert(
        String::from("*.{lua,fcgi,nse,pd_lua,rbxs,wlua}"),
        String::from("Lua"),
    );
    result.insert(String::from("*.{m,moo}"), String::from("Mercury"));
    result.insert(String::from("*.js"), String::from("Node"));
    result.insert(
        String::from("*.{ml,eliom,eliomi,ml4,mli,mll,mly}"),
        String::from("OCaml"),
    );
    result.insert(String::from("*.{m,h}"), String::from("Objective-C"));
    result.insert(String::from("*.opa"), String::from("Opa"));
    result.insert(
        String::from("*.{pl,al,cgi,fcgi,perl,ph,plx,pm,pod,psgi,t}"),
        String::from("Perl"),
    );
    result.insert(String::from("*.pde"), String::from("Processing"));
    result.insert(String::from("*.purs"), String::from("PureScript"));
    result.insert(
        String::from("*.{py,bzl,cgi,fcgi,gyp,lmi,pyde,pyp,pyt,pyw,rpy,tac,wsgi,xpy}"),
        String::from("Python"),
    );
    result.insert(String::from("*.{r,rd,rsx}"), String::from("R"));
    result.insert(
        String::from("*.{rkt,rktd,rktl,scrbl}"),
        String::from("Racket"),
    );
    result.insert(String::from("*.{rb,builder,fcgi,gemspec,god,irbrc,jbuilder,mspec,pluginspec,podspec,rabl,rake,rbuild,rbw,rbx,ru,ruby,thor,watchr}"), String::from("Ruby"));
    result.insert(String::from("*.{rs,rs.in}"), String::from("Rust"));
    result.insert(String::from("*.sass"), String::from("Sass"));
    result.insert(String::from("*.{scala,sbt,sc}"), String::from("Scala"));
    result.insert(
        String::from("*.{scm,sld,sls,sps,ss}"),
        String::from("Scheme"),
    );
    result.insert(String::from("*.{st,cs}"), String::from("Smalltalk"));
    result.insert(String::from("*.swift"), String::from("Swift"));
    result.insert(
        String::from("*.{tex,aux,bbx,bib,cbx,cls,dtx,ins,lbx,ltx,mkii,mkiv,mkvi,sty,toc}"),
        String::from("TeX"),
    );
    result.insert(String::from("*.unity"), String::from("Unity"));
    result.insert(
        String::from("*.{xojo_code,xojo_menu,xojo_report,xojo_script,xojo_toolbar,xojo_window}"),
        String::from("Xojo"),
    );
    result.insert(String::from("*.zep"), String::from("Zephir"));
    result
}

#[cfg(test)]
mod tests {
    use std::io::Error;

    #[test]
    fn finds_single_candidate() -> Result<(), Error> {
        let dir = tempfile::tempdir()?;
        let dir_path = dir.into_path();
        let file_path = dir_path.join("Cargo.toml");
        std::fs::write(&file_path, "")?;
        let result = super::get_language_candidates(&dir_path)?;
        let expected: Vec<String> = vec![String::from("Rust")];
        assert_eq!(result, expected);
        Ok(())
    }

    #[test]
    fn finds_multiple_candidates() -> Result<(), Error> {
        let dir = tempfile::tempdir()?;
        let dir_path = dir.into_path();
        let file1_path = dir_path.join("package.json");
        let file2_path = dir_path.join("pubspec.yaml");
        let file3_path = dir_path.join("go.sum");
        std::fs::write(&file1_path, "")?;
        std::fs::write(&file2_path, "")?;
        std::fs::write(&file3_path, "")?;
        let result = super::get_language_candidates(&dir_path)?;
        let expected: Vec<String> = vec![
            String::from("Dart"),
            String::from("Go"),
            String::from("Node"),
        ];
        assert_eq!(result, expected);
        Ok(())
    }
}