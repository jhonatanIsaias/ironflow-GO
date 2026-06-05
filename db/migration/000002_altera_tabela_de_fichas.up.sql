ALTER TABLE treino.fit_ficha_treino
DROP COLUMN IF EXISTS fit_nr_meta_repeticoes;

ALTER TABLE treino.fit_ficha_treino
ADD COLUMN fit_tx_meta_repeticoes varchar(50) NOT NULL,
ADD COLUMN fit_nr_grupo INTEGER NULL;

ALTER TABLE treino.sex_serie_executada
DROP COLUMN IF EXISTS sex_nr_repeticoes_realizadas;

ALTER TABLE treino.sex_serie_executada
ADD COLUMN sex_tx_repeticoes_executadas varchar(50) NOT NULL;