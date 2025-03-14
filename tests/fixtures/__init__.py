from pytest import fixture


@fixture
def sample_urls_fixtures(faker):
    random_urls = [faker.url() for url in range(10)]
    return random_urls
