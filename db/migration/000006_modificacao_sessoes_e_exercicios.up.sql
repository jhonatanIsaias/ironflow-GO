
ALTER TABLE treino.set_sessao_treino 
ADD COLUMN set_dt_data_fim DATE;

ALTER TABLE treino.set_sessao_treino 
ADD COLUMN set_tm_hora_fim TIME;

ALTER TABLE treino.fit_ficha_treino 
ADD COLUMN exe_tx_tipo_equipamento VARCHAR(50);

ALTER TABLE treino.exe_exercicio 
DROP COLUMN exe_tx_tipo_equipamento;
