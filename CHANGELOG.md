## [1.0.1](https://github.com/juanlugm/concourse-gcs-resource/compare/v1.0.0...v1.0.1) (2025-04-10)


### Bug Fixes

* workflow ([4a0e727](https://github.com/juanlugm/concourse-gcs-resource/commit/4a0e7279d6d8a3fe831c3d21e829880acce75f16))

# 1.0.0 (2025-04-10)


### Bug Fixes

* add debug info to in ([fcba2c7](https://github.com/juanlugm/concourse-gcs-resource/commit/fcba2c76e54c5266a84c04a2ce62bb87f0e67cce))
* Add mission out output ([f9abb78](https://github.com/juanlugm/concourse-gcs-resource/commit/f9abb788450b7714989a566d73bc474c130730c2))
* Add mission out output ([e48f06f](https://github.com/juanlugm/concourse-gcs-resource/commit/e48f06f75a8910a5d94294e586a07dc230bfff81))
* improve version management to allow pinning resources ([333fa73](https://github.com/juanlugm/concourse-gcs-resource/commit/333fa73237746c1f7fc4f45a7936dca7ee103d50))
* Indentation ([286e523](https://github.com/juanlugm/concourse-gcs-resource/commit/286e5234bc1decd4408772f46ae8e6f62abf3cd3))
* object metadata ([4384799](https://github.com/juanlugm/concourse-gcs-resource/commit/4384799f96599c6f8c92b2ea69d045d4ca845066))
* object metadata ([27c0648](https://github.com/juanlugm/concourse-gcs-resource/commit/27c0648c5e2835ca0e8a6645d34f19f58cbcdf4f))
* publish actions ([b934b75](https://github.com/juanlugm/concourse-gcs-resource/commit/b934b75d9ef56e7b3ae0267550366835e5355220))
* select alpine version ([0f229bd](https://github.com/juanlugm/concourse-gcs-resource/commit/0f229bd53dfb9a6f96bc98fe20dd1ced8c88a336))
* Stick to S3 resource conventions. Add readme. ([7ea69eb](https://github.com/juanlugm/concourse-gcs-resource/commit/7ea69eb1bd8775ca52ff6fe8732f579c01d09855))
* take correct params for put ([9122ba2](https://github.com/juanlugm/concourse-gcs-resource/commit/9122ba2a84a7a42b6910b1f3f0e2ef342746eb56))
* typo ([4fbcd53](https://github.com/juanlugm/concourse-gcs-resource/commit/4fbcd534429f3b6b768b2ab1fe820ba0e30c93e7))


### Features

* first implementation ([037aa16](https://github.com/juanlugm/concourse-gcs-resource/commit/037aa16263f1de2dd873164008a332c5655c5f97))

repos:
-   repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v2.3.0
    hooks:
    -   id: check-yaml
    -   id: end-of-file-fixer
    -   id: trailing-whitespace
-   repo: https://github.com/psf/black
    rev: 22.10.0
    hooks:
    -   id: black
