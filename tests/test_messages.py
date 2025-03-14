import unittest

import pytest


class UrlTestCase(unittest.TestCase):

    @pytest.fixture(autouse=True)
    def set_up_fixtures(self, sample_urls_fixtures) -> None:
        self.urls = sample_urls_fixtures

    def test_text_has_a_url(self):
        print(self.urls)

    def test_text_not_have_a_url(self):
        print(self.urls)
